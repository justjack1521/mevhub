package subscriber

import (
	"log/slog"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/player"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
	"time"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type SessionLobbyWriter struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceRepository
	Logger            *slog.Logger
}

func NewSessionLobbyWriter(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, logger *slog.Logger) *SessionLobbyWriter {
	var subscriber = &SessionLobbyWriter{EventPublisher: publisher, SessionRepository: sessions, Logger: logger}
	publisher.Subscribe(subscriber, game.ParticipantDeletedEvent{}, player.DisconnectedEvent{}, player.ConnectedEvent{}, game.PlayerTimedOutEvent{})
	return subscriber
}

func (s *SessionLobbyWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.ParticipantDeletedEvent:
		if err := s.HandleGameParticipantDelete(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("session writer failed to handle event")
		}
	case player.DisconnectedEvent:
		if err := s.HandlePlayerDisconnected(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("session writer failed to handle event")
		}
	case player.ConnectedEvent:
		if err := s.HandlePlayerConnected(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("session writer failed to handle event")
		}
	case game.PlayerTimedOutEvent:
		if err := s.HandlePlayerTimedOut(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("session writer failed to handle event")
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

	// Out-of-order guard: ignore a disconnect that predates the most recently
	// applied connect/disconnect (the two streams can be delivered reordered).
	if event.Timestamp() == 0 {
		// A zero timestamp (older gateway build) silently disables the
		// ordering guard for this session — surface it.
		s.Logger.With("user.id", instance.UserID.String()).Warn("disconnect event carries no timestamp; out-of-order protection inactive")
	}
	if instance.ApplyConnEvent(event.Timestamp()) {
		return nil
	}

	if !uuid.Equal(instance.GameID, uuid.Nil) {
		instance.DisconnectSessionID = event.SessionID()
		if instance.DisconnectedAt = event.Timestamp(); instance.DisconnectedAt == 0 {
			instance.DisconnectedAt = time.Now().UTC().UnixMilli()
		}
		if err := s.SessionRepository.UpdateConnectionState(event.Context(), instance); err != nil {
			return err
		}
		s.EventPublisher.Notify(game.NewPlayerDisconnectedEvent(event.Context(), instance.GameID, instance.UserID, instance.PlayerID))
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

	// Out-of-order guard: ignore a connect that predates the most recently
	// applied connect/disconnect (e.g. a late connect arriving after a newer
	// disconnect), which would otherwise wrongly resume a since-dropped player.
	if event.Timestamp() == 0 {
		s.Logger.With("user.id", instance.UserID.String()).Warn("connect event carries no timestamp; out-of-order protection inactive")
	}
	if instance.ApplyConnEvent(event.Timestamp()) {
		return nil
	}

	if uuid.Equal(instance.GameID, uuid.Nil) || uuid.Equal(instance.DisconnectSessionID, uuid.Nil) {
		// A different gateway session connecting while this session still
		// points at a lobby is an app relaunch: tear the stale session down so
		// the client is not locked out of joining until the key expires.
		if !uuid.Equal(instance.LobbyID, uuid.Nil) &&
			!uuid.Equal(instance.CurrentSessionID, uuid.Nil) &&
			!uuid.Equal(instance.CurrentSessionID, event.SessionID()) &&
			uuid.Equal(instance.GameID, uuid.Nil) {
			if err := s.SessionRepository.Delete(event.Context(), instance); err != nil {
				return err
			}
			s.EventPublisher.Notify(session.NewInstanceDeletedEvent(event.Context(), instance.UserID, instance.PlayerID, instance.LobbyID, instance.GameID, instance.PartySlot, instance.DeckIndex))
			return nil
		}
		// Otherwise just persist the advanced watermark and session identity.
		instance.CurrentSessionID = event.SessionID()
		return s.SessionRepository.UpdateConnectionState(event.Context(), instance)
	}

	if uuid.Equal(instance.DisconnectSessionID, event.SessionID()) {
		// Genuine reconnect within the grace period: the same session has
		// resumed. Clear the disconnect marker and resume the player in-game.
		instance.DisconnectSessionID = uuid.Nil
		instance.DisconnectedAt = 0
		instance.CurrentSessionID = event.SessionID()
		if err := s.SessionRepository.UpdateConnectionState(event.Context(), instance); err != nil {
			return err
		}
		s.EventPublisher.Notify(game.NewPlayerReconnectedEvent(event.Context(), instance.GameID, instance.UserID, instance.PlayerID))
		return nil
	}

	// Different session = app relaunch: kick the stale session.
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
