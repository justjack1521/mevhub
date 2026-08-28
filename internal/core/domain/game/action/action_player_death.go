package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

var (
	ErrFailedPlayerDeath = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to process death for player %s: %w", player, err)
	}
)

// PlayerDeathAction records that a player has fallen in battle. The server
// does not simulate combat, so deaths are client-reported, exactly as enemy HP
// was: the dying client tells the server about its own death. From then on the
// player cannot act until revived, and the death sweep may kick them once the
// game's DeadPlayerKickDuration lapses with no revive.
type PlayerDeathAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	DeathTime  time.Time
}

func NewPlayerDeathAction(instanceID uuid.UUID, playerID uuid.UUID, deathTime time.Time) *PlayerDeathAction {
	return &PlayerDeathAction{InstanceID: instanceID, PlayerID: playerID, DeathTime: deathTime}
}

func (a *PlayerDeathAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerDeath(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerDeath(a.PlayerID, err)
	}

	// Idempotent: a re-reported death (client retry, duplicate delivery) must
	// not reset DeathTime, which would push the kick deadline out forever.
	if player.Dead {
		return nil
	}

	player.Dead = true
	player.DeathTime = a.DeathTime
	// A dead player's queued actions are void. The lock state is deliberately
	// left alone: lock indexes are assigned from the party's locked count, and
	// releasing one mid-turn would let a later lock collide with it.
	player.Actions = nil

	instance.SendChange(NewPlayerDeathChange(instance.InstanceID, a.PlayerID, party.PartyIndex, player.PartySlot))
	return nil
}
