package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type SessionLobbyWriter struct {
	EventPublisher    *mevent.Publisher
	SessionRepository session.InstanceRepository
}

func NewSessionLobbyWriter(publisher *mevent.Publisher, sessions session.InstanceRepository) *SessionLobbyWriter {
	var subscriber = &SessionLobbyWriter{EventPublisher: publisher, SessionRepository: sessions}
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{}, lobby.ParticipantDeletedEvent{})
	return subscriber
}

func (s *SessionLobbyWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.ParticipantCreatedEvent:
		if err := s.HandleParticipantCreate(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantDeletedEvent:
		if err := s.HandleParticipantDelete(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *SessionLobbyWriter) HandleParticipantCreate(event lobby.ParticipantCreatedEvent) error {

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.ClientID())
	if err != nil {
		return err
	}

	instance.LobbyID = event.LobbyID()
	instance.DeckIndex = event.DeckIndex()
	instance.PartySlot = event.SlotIndex()

	if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
		return err
	}

	return nil

}

func (s *SessionLobbyWriter) HandleParticipantDelete(event lobby.ParticipantDeletedEvent) error {

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.ClientID())
	if err != nil {
		return err
	}

	instance.LobbyID = uuid.Nil
	instance.PartySlot = 0

	if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
		return err
	}
	return nil

}
