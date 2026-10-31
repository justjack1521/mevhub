package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

var (
	ErrFailedPlayerRevive = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to revive player %s: %w", player, err)
	}
)

// PlayerReviveAction brings a dead player back into the fight. The reviver may
// be another player or the dead player themselves (a self-revive item), so the
// source is carried separately from the target and is deliberately not
// validated against being alive — the client simulation owns the combat rules
// for who may revive whom.
type PlayerReviveAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	SourceID   uuid.UUID
}

func NewPlayerReviveAction(instanceID uuid.UUID, playerID uuid.UUID, sourceID uuid.UUID) *PlayerReviveAction {
	return &PlayerReviveAction{InstanceID: instanceID, PlayerID: playerID, SourceID: sourceID}
}

func (a *PlayerReviveAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerRevive(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerRevive(a.PlayerID, err)
	}

	// Reviving the living is a no-op, not an error: duplicate revives of the
	// same corpse race each other and the loser did nothing wrong.
	if player.Dead == false {
		return nil
	}

	player.Dead = false
	player.DeathTime = time.Time{}

	// The corpse is gone, so any standing revive claim on it is spent. Clearing
	// here (rather than on death) is what guarantees a fresh death never starts
	// life encumbered by a stale claim.
	player.ReviveClaimSource = uuid.Nil
	player.ReviveClaimExpiry = time.Time{}

	instance.SendChange(NewPlayerReviveChange(instance.InstanceID, a.PlayerID, a.SourceID, party.PartyIndex, player.PartySlot))
	return nil
}
