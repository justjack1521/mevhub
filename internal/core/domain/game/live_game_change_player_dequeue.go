package game

import uuid "github.com/satori/go.uuid"

type PlayerDequeueActionChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerDequeueActionChange(instanceID uuid.UUID, partyIndex int, partySlot int) *PlayerDequeueActionChange {
	return &PlayerDequeueActionChange{InstanceID: instanceID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerDequeueActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierDequeueAction
}
