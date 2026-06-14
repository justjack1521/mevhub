package subscriber

import (
	"fmt"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/player"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type SessionLobbyWriter struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceRepository
}

func NewSessionLobbyWriter(publisher *mevent.Publisher, sessions port.SessionInstanceRepository) *SessionLobbyWriter {
	var subscriber = &SessionLobbyWriter{EventPublisher: publisher, SessionRepository: sessions}
	publisher.Subscribe(subscriber, game.ParticipantDeletedEvent{}, player.DisconnectedEvent{}, player.ConnectedEvent{}, game.PlayerTimedOutEvent{})
	return subscriber
}

func (s *SessionLobbyWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.ParticipantDeletedEvent:
		if err := s.HandleGameParticipantDelete(actual); err != nil {
			fmt.Println(err)
		}
	case player.DisconnectedEvent:
		if err := s.HandlePlayerDisconnected(actual); err != nil {
			fmt.Println(err)
		}
	case player.ConnectedEvent:
		if err := s.HandlePlayerConnected(actual); err != nil {
			fmt.Println(err)
		}
	case game.PlayerTimedOutEvent:
		if err := s.HandlePlayerTimedOut(actual); err != nil {
			fmt.Println(err)
		}
	}
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

	if instance.GameID != event.GameID() {
		return nil
	}

	instance.GameID = uuid.Nil

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

	if !uuid.Equal(instance.GameID, uuid.Nil) {
		instance.DisconnectSessionID = event.SessionID()
		if err := s.SessionRepository.Update(event.Context(), instance); err != nil {
			return err
		}
		s.EventPublisher.Notify(game.NewPlayerDisconnectedEvent(event.Context(), instance.GameID, instance.LobbyID, instance.UserID, instance.PlayerID))
		return nil
	}

	if err := s.SessionRepository.Delete(event.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.UserID, instance.PlayerID, instance.LobbyID, instance.GameID, instance.PartySlot, instance.DeckIndex))

	return nil

}

func (s *SessionLobbyWriter) HandlePlayerConnected(event player.ConnectedEvent) error {

	exists, err := s.SessionRepository.Exists(event.Context(), event.UserID())
	if err != nil || !exists {
		return err
	}

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if uuid.Equal(instance.GameID, uuid.Nil) || uuid.Equal(instance.DisconnectSessionID, uuid.Nil) {
		return nil
	}

	if uuid.Equal(instance.DisconnectSessionID, event.SessionID()) {
		return nil
	}

	if err := s.SessionRepository.Delete(event.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.UserID, instance.PlayerID, instance.LobbyID, instance.GameID, instance.PartySlot, instance.DeckIndex))

	return nil

}

func (s *SessionLobbyWriter) HandlePlayerTimedOut(event game.PlayerTimedOutEvent) error {

	exists, err := s.SessionRepository.Exists(event.Context(), event.UserID())
	if err != nil || !exists {
		return err
	}

	instance, err := s.SessionRepository.QueryByID(event.Context(), event.UserID())
	if err != nil {
		return err
	}

	if !uuid.Equal(instance.GameID, event.GameID()) {
		return nil
	}

	if err := s.SessionRepository.Delete(event.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.UserID, instance.PlayerID, instance.LobbyID, instance.GameID, instance.PartySlot, instance.DeckIndex))

	return nil

}
