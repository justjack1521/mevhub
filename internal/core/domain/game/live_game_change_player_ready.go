package game

import uuid "github.com/satori/go.uuid"

type PlayerReadyChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerReadyChange(id uuid.UUID, party int, slot int) *PlayerReadyChange {
	return &PlayerReadyChange{InstanceID: id, PartyIndex: party, PartySlot: slot}
}

func (c PlayerReadyChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierPlayerReady
}
