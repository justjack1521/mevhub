package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type LobbyInstanceWriter struct {
	EventPublisher          *mevent.Publisher
	LobbyInstanceRepository port.LobbyInstanceRepository
	ParticipantRepository   port.LobbyParticipantRepository
}

func NewLobbyInstanceWriter(publisher *mevent.Publisher, lobbies port.LobbyInstanceRepository, participants port.LobbyParticipantRepository) *LobbyInstanceWriter {
	var service = &LobbyInstanceWriter{EventPublisher: publisher, LobbyInstanceRepository: lobbies, ParticipantRepository: participants}
	publisher.Subscribe(service, session.InstanceDeletedEvent{})
	return service
}

func (s *LobbyInstanceWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case session.InstanceDeletedEvent:
		if err := s.HandleSessionDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyInstanceWriter) HandleSessionDeleted(evt session.InstanceDeletedEvent) error {

	if uuid.Equal(evt.LobbyID(), uuid.Nil) {
		return nil
	}

	instance, err := s.LobbyInstanceRepository.QueryByID(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	if instance.HostPlayerID != evt.PlayerID() {
		participant, err := s.ParticipantRepository.QueryParticipantForLobby(evt.Context(), evt.LobbyID(), evt.PartySlot())
		if err != nil {
			return err
		}
		if err := s.ParticipantRepository.Delete(evt.Context(), participant); err != nil {
			return err
		}
		s.EventPublisher.Notify(lobby.NewParticipantDeletedEvent(evt.Context(), participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))
		return nil
	}

	if err := s.LobbyInstanceRepository.Delete(evt.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(lobby.NewInstanceDeletedEvent(evt.Context(), instance.SysID, instance.QuestID))

	return nil

}
