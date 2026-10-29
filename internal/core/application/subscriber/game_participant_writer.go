package subscriber

import (
	"log/slog"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
)

type GameParticipantWriter struct {
	EventPublisher             *mevent.Publisher
	LobbyParticipantRepository port.LobbyParticipantRepository
	GameParticipantRepository  port.GameParticipantRepository
	LoadoutRepository          port.GamePlayerLoadoutWriteRepository
	Logger                     *slog.Logger
}

func NewGameParticipantWriter(publisher *mevent.Publisher, source port.LobbyParticipantRepository, target port.GameParticipantRepository, loadouts port.GamePlayerLoadoutWriteRepository, logger *slog.Logger) *GameParticipantWriter {
	var service = &GameParticipantWriter{EventPublisher: publisher, LobbyParticipantRepository: source, GameParticipantRepository: target, LoadoutRepository: loadouts, Logger: logger}
	publisher.Subscribe(service, game.PartyCreatedEvent{}, game.PartyDeletedEvent{})
	return service
}

func (s *GameParticipantWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.PartyCreatedEvent:
		if err := s.HandlePartyCreated(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("game participant writer failed to handle event")
		}
	case game.PartyDeletedEvent:
		if err := s.HandlePartyDeleted(actual); err != nil {
			s.Logger.With("event", actual.Name(), "error", err.Error()).Error("game participant writer failed to handle event")
		}
	}
}

func (s *GameParticipantWriter) HandlePartyDeleted(evt game.PartyDeletedEvent) error {
	participants, err := s.GameParticipantRepository.QueryAll(evt.Context(), evt.PartyID())
	if err != nil {
		return err
	}
	for _, participant := range participants {
		s.EventPublisher.Notify(game.NewParticipantDeletedEvent(evt.Context(), evt.GameID(), evt.PartyID(), participant.UserID, participant.PlayerSlot))
	}
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

		if !participant.HasPlayer() {
			continue
		}

		var result = &game.Participant{
			UserID:     participant.UserID,
			PlayerID:   participant.PlayerID,
			PlayerSlot: participant.PlayerSlot,
			DeckIndex:  participant.DeckIndex,
			BotControl: participant.BotControl,
		}

		// Drop any loadout cached by an earlier game so this one fetches the deck as it is now.
		// Session-end eviction alone can miss: a lost SessionEnd, or a teammate's GetGame after
		// the player left, leaves the old deck cached for the full TTL. Evicted before the row is
		// written: clients are told the lobby started before this runs, and GetGame can only
		// reach this player's loadout once the row exists.
		if err := s.LoadoutRepository.Delete(evt.Context(), result.PlayerID, result.DeckIndex); err != nil {
			s.Logger.With("player.id", result.PlayerID.String(), "deck.index", result.DeckIndex, "error", err.Error()).Error("game participant writer failed to evict cached loadout")
		}

		if err := s.GameParticipantRepository.Create(evt.Context(), evt.PartyID(), result); err != nil {
			return err
		}

		s.EventPublisher.Notify(game.NewParticipantCreatedEvent(evt.Context(), evt.GameID(), evt.PartyID(), participant.UserID, participant.PlayerSlot))

	}

	if err := s.LobbyParticipantRepository.DeleteAllForLobby(evt.Context(), evt.PartyID()); err != nil {
		return err
	}

	return nil

}
