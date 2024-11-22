package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
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
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{}, lobby.ParticipantDeletedEvent{}, game.ParticipantCreatedEvent{}, game.ParticipantDeletedEvent{}, player.DisconnectedEvent{})
	return subscriber
}

func (s *SessionLobbyWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.ParticipantCreatedEvent:
		if err := s.HandleLobbyParticipantCreate(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantDeletedEvent:
		if err := s.HandleLobbyParticipantDelete(actual); err != nil {
			fmt.Println(err)
		}
	case game.ParticipantCreatedEvent:
		if err := s.HandleGameParticipantCreate(actual); err != nil {
			fmt.Println(err)
		}
	case game.ParticipantDeletedEvent:
		if err := s.HandleGameParticipantDelete(actual); err != nil {
			fmt.Println(err)
		}
	case player.DisconnectedEvent:
		if err := s.HandlePlayerDisconnected(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *SessionLobbyWriter) HandleGameParticipantCreate(event game.ParticipantCreatedEvent) error {

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

	if instance.LobbyID != event.PartyID() {
		return nil
	}

	instance.GameID = event.GameID()
	if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
		return err
	}

	return nil

}

func (s *SessionLobbyWriter) HandleGameParticipantDelete(event game.ParticipantDeletedEvent) error {

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

	if instance.LobbyID != event.PartyID() {
		return nil
	}

	instance.GameID = uuid.Nil

	if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
		return err
	}

	return nil

}

func (s *SessionLobbyWriter) HandleLobbyParticipantCreate(event lobby.ParticipantCreatedEvent) error {

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

	if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
		return err
	}

	return nil

}

func (s *SessionLobbyWriter) HandleLobbyParticipantDelete(event lobby.ParticipantDeletedEvent) error {

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

	s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.LobbyID, instance.GameID, instance.UserID, instance.PlayerID))

	return nil

}
