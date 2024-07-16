package subscriber

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/consumer"
	"mevhub/internal/domain/game"
)

type GameChannelEventNotifier struct {
	EventPublisher *mevent.Publisher
}

func NewGameChannelEventNotifier(publisher *mevent.Publisher) *GameChannelEventNotifier {
	var subscriber = &GameChannelEventNotifier{EventPublisher: publisher}
	publisher.Subscribe(subscriber, game.InstanceReadyEvent{})
	return subscriber
}

func (s *GameChannelEventNotifier) Notify(event mevent.Event) {

}

func (s *GameChannelEventNotifier) publish(ctx context.Context, t protomulti.MultiGameNotificationType, id uuid.UUID, n ClientNotification) error {
	bytes, err := n.MarshallBinary()
	if err != nil {
		return err
	}
	var message = consumer.NewGameClientNotificationEvent(ctx, t, id, bytes)
	s.EventPublisher.Notify(message)
	return nil
}
