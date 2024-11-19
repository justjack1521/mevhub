package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/player"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type SessionLobbyWriter struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceRepository
}

func NewSessionLobbyWriter(publisher *mevent.Publisher, sessions port.SessionInstanceRepository) *SessionLobbyWriter {
	var subscriber = &SessionLobbyWriter{EventPublisher: publisher, SessionRepository: sessions}
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{}, lobby.ParticipantDeletedEvent{}, player.DisconnectedEvent{})
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
	case player.DisconnectedEvent:
		if err := s.HandlePlayerDisconnected(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *SessionLobbyWriter) HandlePlayerDisconnected(event player.DisconnectedEvent) error {

	exists, err := s.SessionRepository.Exists(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if exists == false {
		return nil
	}

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if err := s.SessionRepository.Delete(event.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.LobbyID, instance.UserID, instance.PlayerID))

	return nil

}

func (s *SessionLobbyWriter) HandleParticipantCreate(event lobby.ParticipantCreatedEvent) error {

	exists, err := s.SessionRepository.Exists(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if exists == false {
		return nil
	}

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.UserID())
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

	exists, err := s.SessionRepository.Exists(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if exists == false {
		return nil
	}

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.UserID())
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
