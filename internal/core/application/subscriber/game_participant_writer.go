package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/factory"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GameParticipantWriter struct {
	EventPublisher             *mevent.Publisher
	LobbyParticipantRepository port.LobbyParticipantReadRepository
	GameParticipantFactory     *factory.PlayerParticipantFactory
	GameParticipantRepository  port.PlayerParticipantWriteRepository
}

func NewGameParticipantWriter(publisher *mevent.Publisher, participants port.LobbyParticipantReadRepository, factory *factory.PlayerParticipantFactory, players port.PlayerParticipantWriteRepository) *GameParticipantWriter {
	var subscriber = &GameParticipantWriter{
		EventPublisher:             publisher,
		LobbyParticipantRepository: participants,
		GameParticipantFactory:     factory,
		GameParticipantRepository:  players,
	}
	publisher.Subscribe(subscriber, game.InstanceRegisteredEvent{})
	return subscriber
}

func (s *GameParticipantWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceRegisteredEvent:
		if err := s.HandleCreate(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GameParticipantWriter) HandleCreate(event game.InstanceRegisteredEvent) error {

	participants, err := s.LobbyParticipantRepository.QueryAllForLobby(event.Context(), event.InstanceID())
	if err != nil {
		return err
	}

	for _, participant := range participants {

		if uuid.Equal(participant.PlayerID, uuid.Nil) {
			continue
		}

		player, err := s.GameParticipantFactory.Create(event.Context(), participant)
		if err != nil {
			return err
		}

		if err := s.GameParticipantRepository.Create(event.Context(), event.InstanceID(), participant.PlayerSlot, player); err != nil {
			return err
		}

		s.EventPublisher.Notify(game.NewParticipantCreatedEvent(event.Context(), event.InstanceID(), participant.PlayerSlot))

	}

	s.EventPublisher.Notify(game.NewInstanceReadyEvent(event.Context(), event.InstanceID()))

	return nil
}
