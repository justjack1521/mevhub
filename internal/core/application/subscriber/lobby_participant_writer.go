package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyParticipantWriter struct {
	EventPublisher        *mevent.Publisher
	ParticipantRepository port.LobbyParticipantRepository
}

func NewLobbyParticipantWriter(publisher *mevent.Publisher, participants port.LobbyParticipantRepository) *LobbyParticipantWriter {
	var service = &LobbyParticipantWriter{EventPublisher: publisher, ParticipantRepository: participants}
	publisher.Subscribe(service, lobby.InstanceDeletedEvent{})
	return service
}

func (s *LobbyParticipantWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceDeletedEvent:
		if err := s.HandleLobbyInstanceDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyParticipantWriter) HandleLobbyInstanceDeleted(evt lobby.InstanceDeletedEvent) error {

	participants, err := s.ParticipantRepository.QueryAllForLobby(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if err := s.ParticipantRepository.Delete(evt.Context(), participant); err != nil {
			return err
		}
		s.EventPublisher.Notify(lobby.NewParticipantDeletedEvent(evt.Context(), participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))
	}

	return nil

}
