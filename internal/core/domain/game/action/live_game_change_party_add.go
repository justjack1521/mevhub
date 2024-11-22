package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PartyAddChange struct {
	PartyID   uuid.UUID
	PartySlot int
}

func NewPartyAddChange(partyID uuid.UUID, partySlot int) *PartyAddChange {
	return &PartyAddChange{PartyID: partyID, PartySlot: partySlot}
}

func (c PartyAddChange) Identifier() game.ChangeIdentifier {
	return game.ChaneIdentifierPartyAdd
}
