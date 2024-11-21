package subscriber

import (
	"errors"
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GamePartyWriter struct {
	EventPublisher         *mevent.Publisher
	GameInstanceRepository port.GameInstanceReadRepository
	LobbySummaryRepository port.LobbySummaryReadRepository
	PartyRepository        port.GamePartyWriteRepository
}

func NewGamePartyWriter(publisher *mevent.Publisher, games port.GameInstanceReadRepository, lobbies port.LobbySummaryReadRepository, parties port.GamePartyWriteRepository) *GamePartyWriter {
	var service = &GamePartyWriter{EventPublisher: publisher, GameInstanceRepository: games, LobbySummaryRepository: lobbies, PartyRepository: parties}
	publisher.Subscribe(service, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{})
	return service
}

func (s *GamePartyWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		if err := s.HandleInstanceCreated(actual); err != nil {
			fmt.Println(err)
		}
	case game.InstanceDeletedEvent:
		if err := s.HandleInstanceDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GamePartyWriter) HandleInstanceDeleted(evt game.InstanceDeletedEvent) error {
	if err := s.PartyRepository.DeleteAll(evt.Context(), evt.InstanceID()); err != nil {
		return err
	}
	return nil
}

func (s *GamePartyWriter) HandleInstanceCreated(evt game.InstanceCreatedEvent) error {

	parent, err := s.GameInstanceRepository.Get(evt.Context(), evt.InstanceID())
	if err != nil {
		return err
	}

	if len(parent.LobbyIDs) == 0 {
		return errors.New("invalid number of lobbies in game")
	}

	for index, value := range parent.LobbyIDs {

		instance, err := s.LobbySummaryRepository.Query(evt.Context(), value)
		if err != nil {
			return err
		}

		result := &game.Party{
			SysID:     instance.InstanceID,
			PartyID:   instance.PartyID,
			Index:     index,
			PartyName: instance.LobbyComment,
		}

		if err := s.PartyRepository.Create(evt.Context(), evt.InstanceID(), result); err != nil {
			return err
		}

		s.EventPublisher.Notify(game.NewPartyCreatedEvent(evt.Context(), result.SysID, parent.SysID))

	}

	return nil

}
