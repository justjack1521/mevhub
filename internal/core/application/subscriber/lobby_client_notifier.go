package subscriber

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/application/consumer"
	"strings"
)

type LobbyClientNotifier struct {
	client *redis.Client
}

func NewLobbyClientNotifier(publisher *mevent.Publisher, redis *redis.Client) *LobbyClientNotifier {
	var subscriber = &LobbyClientNotifier{client: redis}
	publisher.Subscribe(subscriber, consumer.LobbyClientNotificationEvent{}, consumer.GameClientNotificationEvent{})
	return subscriber
}

func (s *LobbyClientNotifier) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case consumer.LobbyClientNotificationEvent:
		s.NotifyLobby(actual)
	case consumer.GameClientNotificationEvent:
		s.NotifyGame(actual)
	}
}

func (s *LobbyClientNotifier) NotifyLobby(event consumer.LobbyClientNotificationEvent) {
	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(event.Operation()),
		Data:    event.Data(),
	}
	bytes, err := notification.MarshallBinary()
	if err != nil {
		fmt.Println(err)
		return
	}
	var key = strings.Join([]string{LobbyChannelPrefix, event.LobbyID().String()}, LobbyChannelSeparator)
	s.client.Publish(event.Context(), key, bytes)
}

func (s *LobbyClientNotifier) NotifyGame(event consumer.GameClientNotificationEvent) {
	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(event.Operation()),
		Data:    event.Data(),
	}
	bytes, err := notification.MarshallBinary()
	if err != nil {
		fmt.Println(err)
		return
	}
	var key = strings.Join([]string{LobbyChannelPrefix, event.GameID().String()}, LobbyChannelSeparator)
	s.client.Publish(event.Context(), key, bytes)
}
