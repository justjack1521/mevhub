package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrFailedEnqueueAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to enqueue action for player %s: %w", player, err)
	}
)

type PlayerEnqueueAction struct {
	InstanceID uuid.UUID
	PartyID    uuid.UUID
	PlayerID   uuid.UUID
	Target     int
	ActionType PlayerActionType
	SlotIndex  int
	ElementID  uuid.UUID
}

func NewPlayerEnqueueAction(instanceID, partyID, playerID uuid.UUID, target int, actionType PlayerActionType, slotIndex int, elementID uuid.UUID) *PlayerEnqueueAction {
	return &PlayerEnqueueAction{InstanceID: instanceID, PartyID: partyID, PlayerID: playerID, Target: target, ActionType: actionType, SlotIndex: slotIndex, ElementID: elementID}
}

func (a *PlayerEnqueueAction) Perform(game *LiveGameInstance) error {

	party, err := game.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedEnqueueAction(a.PlayerID, err)
	}

	var action = &PlayerAction{
		Target:     a.Target,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		ElementID:  a.ElementID,
	}

	if err := player.EnqueueAction(action); err != nil {
		return ErrFailedEnqueueAction(a.PlayerID, err)
	}

	game.SendChange(NewPlayerEnqueueActionChange(game.InstanceID, party.PartyIndex, player.PartySlot, a.ActionType, a.SlotIndex, a.Target, a.ElementID))

	return nil

}
