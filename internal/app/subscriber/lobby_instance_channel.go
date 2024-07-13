package subscriber

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/justjack1521/mevium/pkg/rabbitmv"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/app/consumer"
	"mevhub/internal/domain/lobby"
	"strings"
)

const LobbyChannelPrefix string = "lobby_channel"
const LobbyChannelSeparator string = ":"

type LobbyNotificationChanneler struct {
	client     *redis.Client
	repository lobby.NotificationListenerRepository
	publisher  *rabbitmv.StandardPublisher
	channels   map[uuid.UUID]*LobbyInstanceNotificationChannel
}

func NewLobbyNotificationChanneler(publisher *mevent.Publisher, client *redis.Client, conn *rabbitmq.Conn, listeners lobby.NotificationListenerRepository) *LobbyNotificationChanneler {
	var manager = &LobbyNotificationChanneler{
		client:     client,
		channels:   make(map[uuid.UUID]*LobbyInstanceNotificationChannel),
		repository: listeners,
		publisher:  rabbitmv.NewClientPublisher(conn, rabbitmq.WithPublisherOptionsLogger(logrus.New())),
	}
	var channels = []mevent.Event{
		mevent.ApplicationStartEvent{},
		mevent.ApplicationShutdownEvent{},
		lobby.InstanceCreatedEvent{},
		lobby.InstanceDeletedEvent{},
		lobby.WatcherAddedEvent{},
		lobby.ParticipantCreatedEvent{},
		lobby.ParticipantDeletedEvent{},
		consumer.LobbyClientNotificationEvent{},
	}
	publisher.Subscribe(manager, channels...)
	return manager
}

func (s *LobbyNotificationChanneler) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case mevent.ApplicationStartEvent:
		s.Start(actual)
	case lobby.InstanceCreatedEvent:
		s.HandleCreate(actual)
	case lobby.InstanceDeletedEvent:
		s.HandleDelete(actual)
	case lobby.WatcherAddedEvent:
		s.HandleWatcherAdd(actual)
	case lobby.ParticipantCreatedEvent:
		s.HandleParticipantAdd(actual)
	case lobby.ParticipantDeletedEvent:
		s.HandleParticipantDelete(actual)
	case mevent.ApplicationShutdownEvent:
		s.CloseAll(actual)
	}
}

func (s *LobbyNotificationChanneler) Start(event mevent.ApplicationStartEvent) {

}

func (s *LobbyNotificationChanneler) CloseAll(event mevent.ApplicationShutdownEvent) {
	for _, channel := range s.channels {
		channel.close(context.Background())
	}
}

func (s *LobbyNotificationChanneler) HandleCreate(event lobby.InstanceCreatedEvent) {
	var channel = s.NewLobbyInstanceNotificationChannel(event.Context(), event.LobbyID(), s)
	s.channels[event.LobbyID()] = channel
	go channel.run()
}

func (s *LobbyNotificationChanneler) HandleDelete(event lobby.InstanceDeletedEvent) {

	defer func() {
		channel, exists := s.channels[event.LobbyID()]
		if exists == false || channel == nil {
			return
		}
		channel.close(event.Context())
		delete(s.channels, event.LobbyID())
	}()

	listeners, err := s.repository.QueryAllForLobby(event.Context(), event.LobbyID())
	if err != nil {
		return
	}

	if err := s.repository.DeleteAll(event.Context(), event.LobbyID()); err != nil {
		return
	}

	var notification = &protomulti.LobbyCancelNotification{LobbyId: event.LobbyID().String()}

	n, err := notification.MarshallBinary()
	if err != nil {
		return
	}

	var message = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(protomulti.MultiNotificationType_LOBBY_CANCEL),
		Data:    n,
	}

	bytes, err := message.MarshallBinary()
	if err != nil {
		return
	}

	for _, listener := range listeners {
		if err := s.publisher.Publish(event.Context(), bytes, listener.UserID, listener.PlayerID, rabbitmv.ClientNotification); err != nil {
			return
		}
	}

}

func (s *LobbyNotificationChanneler) HandleParticipantAdd(event lobby.ParticipantCreatedEvent) {
	channel, exists := s.channels[event.LobbyID()]
	if exists == false || channel == nil {
		return
	}
	if err := s.repository.Create(event.Context(), event.LobbyID(), event.UserID()); err != nil {
		return
	}
}

func (s *LobbyNotificationChanneler) HandleParticipantDelete(event lobby.ParticipantDeletedEvent) {
	channel, exists := s.channels[event.LobbyID()]
	if exists == false || channel == nil {
		return
	}
	if err := s.repository.Delete(event.Context(), event.LobbyID(), event.UserID()); err != nil {
		return
	}
}

func (s *LobbyNotificationChanneler) HandleWatcherAdd(event lobby.WatcherAddedEvent) {
	channel, exists := s.channels[event.LobbyID()]
	if exists == false || channel == nil {
		return
	}
	if err := s.repository.Create(event.Context(), event.LobbyID(), event.UserID()); err != nil {
		return
	}
}

func (s *LobbyNotificationChanneler) NewLobbyInstanceNotificationChannel(ctx context.Context, instance uuid.UUID, manager *LobbyNotificationChanneler) *LobbyInstanceNotificationChannel {
	var channel = &LobbyInstanceNotificationChannel{
		LobbyID: instance,
		manager: manager,
		channel: s.client.Subscribe(ctx, strings.Join([]string{LobbyChannelPrefix, instance.String()}, LobbyChannelSeparator)),
	}
	return channel
}

type LobbyInstanceNotificationChannel struct {
	LobbyID uuid.UUID
	channel *redis.PubSub
	manager *LobbyNotificationChanneler
}

func (c *LobbyInstanceNotificationChannel) Publish(ctx context.Context, message []byte) error {

	listeners, err := c.manager.repository.QueryAllForLobby(ctx, c.LobbyID)

	if err != nil {
		return err
	}

	if len(listeners) == 0 {
		c.close(ctx)
	}

	for _, listener := range listeners {
		if err := c.manager.publisher.Publish(ctx, message, listener.UserID, listener.PlayerID, rabbitmv.ClientNotification); err != nil {
			return err
		}
	}
	return nil
}

func (c *LobbyInstanceNotificationChannel) run() {
	channel := c.channel.Channel()
	for message := range channel {
		if err := c.Publish(context.Background(), []byte(message.Payload)); err != nil {
			fmt.Println("Error Publishing Message To Listener: ", err)
		}
	}
}

func (c *LobbyInstanceNotificationChannel) close(ctx context.Context) {
	if err := c.channel.Unsubscribe(ctx, c.LobbyID.String()); err != nil {
		fmt.Println(err)
	}
}
