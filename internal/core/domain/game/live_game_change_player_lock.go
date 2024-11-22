package game

import uuid "github.com/satori/go.uuid"

type PlayerLockActionChange struct {
	InstanceID      uuid.UUID
	PartyIndex      int
	PartySlot       int
	ActionLockIndex int
}

func NewPlayerLockActionChange(id uuid.UUID, party int, slot int, index int) *PlayerLockActionChange {
	return &PlayerLockActionChange{InstanceID: id, PartyIndex: party, PartySlot: slot, ActionLockIndex: index}
}

func (c PlayerLockActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierLockAction
}
