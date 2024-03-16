package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/domain/lobby"
)

type LobbySearchWriter struct {
	EventPublisher   *mevent.Publisher
	SearchRepository lobby.SearchWriteRepository
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

	var search = lobby.SearchEntry{
		InstanceID:         event.LobbyID(),
		ModeIdentifier:     event.Mode(),
		Level:              event.Level(),
		MinimumPlayerLevel: event.MinLevel(),
		Categories:         event.Categories(),
	}

	if err := s.SearchRepository.Create(event.Context(), search); err != nil {
		return err
	}

	return nil

}
