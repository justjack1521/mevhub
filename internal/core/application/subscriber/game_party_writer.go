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
	fmt.Println("Receive event")
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		fmt.Println("Receive Instance created")
		if err := s.HandleInstanceCreated(actual); err != nil {
			fmt.Println("Instance create failed")
			fmt.Println(err)
		}
		fmt.Println("Instance create success")
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
		fmt.Println("Failed get instance")
		return err
	}

	fmt.Println("Got instance")

	if len(parent.LobbyIDs) == 0 {
		fmt.Println("no lobbies")
		return errors.New("invalid number of lobbies in game")
	}

	for index, value := range parent.LobbyIDs {
		fmt.Println("lobby id", value)
		instance, err := s.LobbySummaryRepository.Query(evt.Context(), value)
		if err != nil {
			fmt.Println("no lobby", value)
			return err
		}
		fmt.Println("got lobby", value)
		result := &game.Party{
			SysID:     instance.InstanceID,
			PartyID:   instance.PartyID,
			Index:     index,
			PartyName: instance.LobbyComment,
		}

		if err := s.PartyRepository.Create(evt.Context(), evt.InstanceID(), result); err != nil {
			fmt.Println("failed create party", value)
			return err
		}
		fmt.Println("create party", value)

		s.EventPublisher.Notify(game.NewPartyCreatedEvent(evt.Context(), result.SysID, parent.SysID, result.Index))

	}

	return nil

}
