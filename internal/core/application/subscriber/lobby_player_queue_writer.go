package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyPlayerQueueWriter struct {
	QueueRepository port.MatchLobbyPlayerQueueWriteRepository
	QuestRepository port.QuestRepository
}

func NewLobbyPlayerQueueWriter(publisher *mevent.Publisher, queues port.MatchLobbyPlayerQueueWriteRepository, quests port.QuestRepository) *LobbyPlayerQueueWriter {
	var subscriber = &LobbyPlayerQueueWriter{QueueRepository: queues, QuestRepository: quests}
	publisher.Subscribe(subscriber, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbyPlayerQueueWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceDeletedEvent:
		if err := s.HandleInstanceDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyPlayerQueueWriter) HandleInstanceDeleted(evt lobby.InstanceDeletedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return nil
	}

	if err := s.QueueRepository.RemoveLobbyFromQueue(evt.Context(), quest.Tier.GameMode.ModeIdentifier, evt.QuestID(), evt.LobbyID()); err != nil {
		return err
	}

	return nil

}
