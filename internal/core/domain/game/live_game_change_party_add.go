package game

import uuid "github.com/satori/go.uuid"

type PartyAddChange struct {
	PartyID   uuid.UUID
	PartySlot int
}

func NewPartyAddChange(partyID uuid.UUID, partySlot int) *PartyAddChange {
	return &PartyAddChange{PartyID: partyID, PartySlot: partySlot}
}

func (c PartyAddChange) Identifier() ChangeIdentifier {
	return ChaneIdentifierPartyAdd
}
