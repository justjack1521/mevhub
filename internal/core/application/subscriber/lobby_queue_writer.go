package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyQueueWriter struct {
	LobbyRepository port.LobbyInstanceReadRepository
	QuestRepository port.QuestRepository
	QueueRepository port.MatchLobbyQueueWriteRepository
}

func NewLobbyQueueWriter(publisher *mevent.Publisher, instance port.LobbyInstanceReadRepository, quests port.QuestRepository, queue port.MatchLobbyQueueWriteRepository) *LobbyQueueWriter {
	var subscriber = &LobbyQueueWriter{LobbyRepository: instance, QuestRepository: quests, QueueRepository: queue}
	publisher.Subscribe(subscriber, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbyQueueWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceDeletedEvent:
		if err := s.HandleLobbyDelete(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyQueueWriter) HandleLobbyDelete(evt lobby.InstanceDeletedEvent) error {

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
