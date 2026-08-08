package action

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

var (
	ErrFailedLockAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s: %w", player, err)
	}
	ErrPlayerUnableToLockAction = errors.New("player unable to dequeue action")
)

type PlayerLockAction struct {
	InstanceID uuid.UUID
	PartyID    uuid.UUID
	PlayerID   uuid.UUID
}

func NewPlayerLockAction(instanceID uuid.UUID, partyID uuid.UUID, playerID uuid.UUID) *PlayerLockAction {
	return &PlayerLockAction{InstanceID: instanceID, PartyID: partyID, PlayerID: playerID}
}

func (a *PlayerLockAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return err
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedLockAction(a.PlayerID, err)
	}

	if player.ActionsLocked {
		return ErrFailedLockAction(a.PlayerID, ErrPlayerUnableToLockAction)
	}

	// Party-scoped, not game-wide: the lock index is consumed as an index
	// into a party-sized queue when the enemy turn is built. A game-wide
	// count with multiple parties produced out-of-range indexes.
	player.ActionLockIndex = party.GetActionLockedPlayerCount()
	player.ActionsLocked = true

	instance.SendChange(NewPlayerLockActionChange(a.InstanceID, party.PartyIndex, player.PartySlot, player.ActionLockIndex))

	return nil

}
