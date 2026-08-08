package subscriber

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/justjack1521/mevrabbit"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"mevhub/internal/core/domain/lobby"
	"strings"
	"sync"
)

const LobbyChannelPrefix string = "lobby_notification_channel"
const LobbyChannelSeparator string = ":"

type LobbyNotificationChanneler struct {
	client     *redis.Client
	repository lobby.NotificationListenerRepository
	publisher  *mevrabbit.StandardPublisher
	logger     *slog.Logger
	// mu guards channels: Open runs on gRPC goroutines while the event
	// handlers run on whichever goroutine published the event.
	mu       sync.Mutex
	channels map[uuid.UUID]*LobbyInstanceNotificationChannel
}

func NewLobbyNotificationChanneler(publisher *mevent.Publisher, client *redis.Client, conn *rabbitmq.Conn, listeners lobby.NotificationListenerRepository, logger *slog.Logger) *LobbyNotificationChanneler {
	var manager = &LobbyNotificationChanneler{
		client:     client,
		channels:   make(map[uuid.UUID]*LobbyInstanceNotificationChannel),
		repository: listeners,
		publisher:  mevrabbit.NewClientPublisher(conn, rabbitmq.WithPublisherOptionsLogger(logrus.New())),
		logger:     logger,
	}
	var channels = []mevent.Event{
		mevent.ApplicationStartEvent{},
		mevent.ApplicationShutdownEvent{},
		lobby.InstanceDeletedEvent{},
		lobby.WatcherAddedEvent{},
		lobby.ParticipantDeletedEvent{},
	}
	publisher.Subscribe(manager, channels...)
	return manager
}

func (s *LobbyNotificationChanneler) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case mevent.ApplicationStartEvent:
		s.Start(actual)
	case lobby.InstanceDeletedEvent:
		s.HandleDelete(actual)
	case lobby.WatcherAddedEvent:
		s.HandleWatcherAdd(actual)
	case lobby.ParticipantDeletedEvent:
		s.HandleParticipantDelete(actual)
	case mevent.ApplicationShutdownEvent:
		s.CloseAll(actual)
	}
}

func (s *LobbyNotificationChanneler) Open(ctx context.Context, id uuid.UUID) {
	var channel = s.NewLobbyInstanceNotificationChannel(ctx, id, s)
	s.mu.Lock()
	previous := s.channels[id]
	s.channels[id] = channel
	s.mu.Unlock()
	if previous != nil {
		previous.close(ctx)
	}
	go channel.run()
}

// remove closes a lobby's channel and drops it from the registry; safe to call
// from any goroutine and for ids that are no longer present.
func (s *LobbyNotificationChanneler) remove(ctx context.Context, id uuid.UUID) {
	s.mu.Lock()
	channel := s.channels[id]
	delete(s.channels, id)
	s.mu.Unlock()
	if channel != nil {
		channel.close(ctx)
	}
}

func (s *LobbyNotificationChanneler) Start(event mevent.ApplicationStartEvent) {
	s.closeAll()
}

func (s *LobbyNotificationChanneler) CloseAll(event mevent.ApplicationShutdownEvent) {
	s.closeAll()
}

func (s *LobbyNotificationChanneler) closeAll() {
	s.mu.Lock()
	channels := make([]*LobbyInstanceNotificationChannel, 0, len(s.channels))
	for _, channel := range s.channels {
		channels = append(channels, channel)
	}
	s.channels = make(map[uuid.UUID]*LobbyInstanceNotificationChannel)
	s.mu.Unlock()
	for _, channel := range channels {
		channel.close(context.Background())
	}
}

func (s *LobbyNotificationChanneler) HandleDelete(event lobby.InstanceDeletedEvent) {

	defer s.remove(event.Context(), event.LobbyID())

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
		Type:    int32(protomulti.MultiLobbyNotificationType_LOBBY_NOTIFY_CANCEL),
		Data:    n,
	}

	bytes, err := message.MarshallBinary()
	if err != nil {
		return
	}

	for _, listener := range listeners {
		if err := s.publisher.Publish(event.Context(), bytes, listener.UserID, listener.PlayerID, mevrabbit.ClientNotification); err != nil {
			return
		}
	}

}

func (s *LobbyNotificationChanneler) HandleParticipantDelete(event lobby.ParticipantDeletedEvent) {
	s.mu.Lock()
	channel, exists := s.channels[event.LobbyID()]
	s.mu.Unlock()
	if exists == false || channel == nil {
		return
	}
	if err := s.repository.DeleteListener(event.Context(), event.LobbyID(), event.UserID()); err != nil {
		return
	}
}

func (s *LobbyNotificationChanneler) HandleWatcherAdd(event lobby.WatcherAddedEvent) {
	s.mu.Lock()
	channel, exists := s.channels[event.LobbyID()]
	s.mu.Unlock()
	if exists == false || channel == nil {
		return
	}
	if err := s.repository.CreateListener(event.Context(), event.LobbyID(), event.UserID(), event.PlayerID()); err != nil {
		return
	}
}

func (s *LobbyNotificationChanneler) NewLobbyInstanceNotificationChannel(ctx context.Context, instance uuid.UUID, manager *LobbyNotificationChanneler) *LobbyInstanceNotificationChannel {
	var channel = &LobbyInstanceNotificationChannel{
		LobbyID: instance,
		manager: manager,
		// The subscription outlives the request that opened the lobby, so it
		// must not be bound to the caller's (cancellable) context.
		channel: s.client.Subscribe(context.Background(), s.Key(instance)),
	}
	return channel
}

func (s *LobbyNotificationChanneler) Key(id uuid.UUID) string {
	return strings.Join([]string{LobbyChannelPrefix, id.String()}, LobbyChannelSeparator)
}
