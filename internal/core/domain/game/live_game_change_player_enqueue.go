package game

import uuid "github.com/satori/go.uuid"

type PlayerEnqueueActionChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
	ActionType PlayerActionType
	SlotIndex  int
	Target     int
	ElementID  uuid.UUID
}

func NewPlayerEnqueueActionChange(instanceID uuid.UUID, partyIndex int, partySlot int, actionType PlayerActionType, slotIndex int, target int, elementID uuid.UUID) *PlayerEnqueueActionChange {
	return &PlayerEnqueueActionChange{InstanceID: instanceID, PartyIndex: partyIndex, PartySlot: partySlot, ActionType: actionType, SlotIndex: slotIndex, Target: target, ElementID: elementID}
}

func (c PlayerEnqueueActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierEnqueueAction
}
