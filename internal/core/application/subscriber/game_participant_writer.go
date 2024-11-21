package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GameParticipantWriter struct {
	EventPublisher             *mevent.Publisher
	LobbyParticipantRepository port.LobbyParticipantReadRepository
	GameParticipantRepository  port.GameParticipantWriteRepository
}

func NewGameParticipantWriter(publisher *mevent.Publisher, source port.LobbyParticipantReadRepository, target port.GameParticipantWriteRepository) *GameParticipantWriter {
	var service = &GameParticipantWriter{EventPublisher: publisher, LobbyParticipantRepository: source, GameParticipantRepository: target}
	publisher.Subscribe(service, game.PartyCreatedEvent{}, game.PartyDeletedEvent{})
	return service
}

func (s *GameParticipantWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.PartyCreatedEvent:
		if err := s.HandlePartyCreated(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GameParticipantWriter) HandlePartyDeleted(evt game.PartyDeletedEvent) error {
	if err := s.GameParticipantRepository.DeleteAll(evt.Context(), evt.PartyID()); err != nil {
		return err
	}
	return nil
}

func (s *GameParticipantWriter) HandlePartyCreated(evt game.PartyCreatedEvent) error {

	participants, err := s.LobbyParticipantRepository.QueryAllForLobby(evt.Context(), evt.PartyID())
	if err != nil {
		return err
	}

	for _, participant := range participants {

		var result = &game.Participant{
			UserID:     participant.UserID,
			PlayerID:   participant.PlayerID,
			PlayerSlot: participant.PlayerSlot,
			DeckIndex:  participant.DeckIndex,
			BotControl: participant.BotControl,
		}

		if err := s.GameParticipantRepository.Create(evt.Context(), evt.PartyID(), result); err != nil {
			return err
		}

		s.EventPublisher.Notify(game.NewParticipantCreatedEvent(evt.Context(), evt.GameID(), evt.PartyID(), participant.UserID, participant.PlayerSlot))

	}

	return nil

}
