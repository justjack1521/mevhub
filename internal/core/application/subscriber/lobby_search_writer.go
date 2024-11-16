package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbySearchWriter struct {
	EventPublisher   *mevent.Publisher
	SearchRepository lobby.SearchWriteRepository
	QuestRepository  port.QuestRepository
}

func NewLobbySearchWriter(publisher *mevent.Publisher, searcher lobby.SearchWriteRepository) *LobbySearchWriter {
	var subscriber = &LobbySearchWriter{EventPublisher: publisher, SearchRepository: searcher}
	publisher.Subscribe(subscriber, lobby.SummaryCreatedEvent{})
	return subscriber
}

func (s *LobbySearchWriter) Notify(event mevent.Event) {
	actual, valid := event.(lobby.SummaryCreatedEvent)
	if valid == false {
		return
	}
	if err := s.Handle(actual); err != nil {
		fmt.Println(err)
	}
}

func (s *LobbySearchWriter) Handle(event lobby.SummaryCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(event.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodSearch {
		return nil
	}

	var search = lobby.SearchEntry{
		InstanceID:         event.LobbyID(),
		ModeIdentifier:     string(quest.Tier.GameMode.ModeIdentifier),
		Level:              event.Level(),
		MinimumPlayerLevel: event.MinLevel(),
		Categories:         event.Categories(),
	}

	if err := s.SearchRepository.Create(event.Context(), search); err != nil {
		return err
	}

	return nil

}
