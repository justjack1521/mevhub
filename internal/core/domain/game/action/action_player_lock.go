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

	party, err := instance.GetParty(a.PartyID)
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

	player.ActionLockIndex = instance.GetActionLockedPlayerCount()
	player.ActionsLocked = true

	instance.SendChange(NewPlayerLockActionChange(a.InstanceID, party.PartyIndex, player.PartySlot, player.ActionLockIndex))

	return nil

}
