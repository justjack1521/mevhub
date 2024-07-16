package subscriber

import (
	"context"
	"fmt"
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
	switch actual := event.(type) {
	case game.InstanceReadyEvent:
		if err := s.HandleReady(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GameChannelEventNotifier) HandleReady(event game.InstanceReadyEvent) error {

	var notification = &protomulti.GameReadyNotification{
		GameId: event.InstanceID().String(),
	}

	return s.publish(event.Context(), protomulti.MultiGameNotificationType_GAME_READY, event.InstanceID(), notification)

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
