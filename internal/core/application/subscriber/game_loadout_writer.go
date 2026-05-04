package subscriber

import (
	"fmt"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
)

type GameLoadoutEvictionSubscriber struct {
	LoadoutRepository port.GamePlayerLoadoutRepository
}

func NewGameLoadoutEvictionSubscriber(publisher *mevent.Publisher, loadouts port.GamePlayerLoadoutRepository) *GameLoadoutEvictionSubscriber {
	var subscriber = &GameLoadoutEvictionSubscriber{LoadoutRepository: loadouts}
	publisher.Subscribe(subscriber, session.InstanceDeletedEvent{})
	return subscriber
}

func (s *GameLoadoutEvictionSubscriber) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case session.InstanceDeletedEvent:
		if err := s.HandleSessionDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GameLoadoutEvictionSubscriber) HandleSessionDeleted(event session.InstanceDeletedEvent) error {
	return s.LoadoutRepository.Delete(event.Context(), event.PlayerID(), event.DeckIndex())
}
