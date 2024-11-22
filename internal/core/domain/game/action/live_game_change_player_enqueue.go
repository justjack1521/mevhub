package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerEnqueueActionChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
	ActionType game.PlayerActionType
	SlotIndex  int
	Target     int
	ElementID  uuid.UUID
}

func NewPlayerEnqueueActionChange(instanceID uuid.UUID, partyIndex int, partySlot int, actionType game.PlayerActionType, slotIndex int, target int, elementID uuid.UUID) *PlayerEnqueueActionChange {
	return &PlayerEnqueueActionChange{InstanceID: instanceID, PartyIndex: partyIndex, PartySlot: partySlot, ActionType: actionType, SlotIndex: slotIndex, Target: target, ElementID: elementID}
}

func (c PlayerEnqueueActionChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierEnqueueAction
}
