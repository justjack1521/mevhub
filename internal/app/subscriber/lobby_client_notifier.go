package subscriber

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/app/consumer"
	"strings"
)

type LobbyClientNotifier struct {
	client *redis.Client
}

func NewLobbyClientNotifier(publisher *mevent.Publisher, redis *redis.Client) *LobbyClientNotifier {
	var subscriber = &LobbyClientNotifier{client: redis}
	publisher.Subscribe(subscriber, consumer.LobbyClientNotificationEvent{})
	return subscriber
}

func (s *LobbyClientNotifier) Notify(event mevent.Event) {
	actual, valid := event.(consumer.LobbyClientNotificationEvent)
	if valid == false {
		return
	}
	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(actual.Operation()),
		Data:    actual.Data(),
	}
	bytes, err := notification.MarshallBinary()
	if err != nil {
		fmt.Println(err)
	}
	var key = strings.Join([]string{LobbyChannelPrefix, actual.LobbyID().String()}, LobbyChannelSeparator)
	s.client.Publish(actual.Context(), key, bytes)
}
