package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbySearchWriter struct {
	QuestRepository  port.QuestRepository
	SearchRepository port.LobbySearchWriteRepository
}

func NewLobbySearchWriter(publisher *mevent.Publisher, quests port.QuestRepository, search port.LobbySearchWriteRepository) *LobbySearchWriter {
	var subscriber = &LobbySearchWriter{QuestRepository: quests, SearchRepository: search}
	publisher.Subscribe(subscriber, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbySearchWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceDeletedEvent:
		if err := s.HandleLobbyDelete(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbySearchWriter) HandleLobbyDelete(evt lobby.InstanceDeletedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodSearch {
		return nil
	}

	var entry = lobby.SearchEntry{
		InstanceID:     evt.LobbyID(),
		ModeIdentifier: string(quest.Tier.GameMode.ModeIdentifier),
		Level:          quest.Tier.StarLevel,
		Categories:     quest.CategoryIDs(),
	}

	if err := s.SearchRepository.Delete(evt.Context(), entry); err != nil {
		return err
	}

	return nil

}
