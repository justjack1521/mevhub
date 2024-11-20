package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GameParticipantWriter struct {
	EventPublisher             *mevent.Publisher
	ParticipantReadRepository  port.LobbyParticipantReadRepository
	ParticipantWriteRepository port.GameParticipantWriteRepository
}

func NewGameParticipantWriter(publisher *mevent.Publisher, source port.LobbyParticipantReadRepository, target port.GameParticipantWriteRepository) *GameParticipantWriter {
	var service = &GameParticipantWriter{EventPublisher: publisher, ParticipantReadRepository: source, ParticipantWriteRepository: target}
	publisher.Subscribe(service, game.PartyCreatedEvent{})
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

func (s *GameParticipantWriter) HandlePartyCreated(evt game.PartyCreatedEvent) error {

	participants, err := s.ParticipantReadRepository.QueryAllForLobby(evt.Context(), evt.PartyID())
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

		if err := s.ParticipantWriteRepository.Create(evt.Context(), evt.PartyID(), result); err != nil {
			return err
		}

	}

	return nil

}
