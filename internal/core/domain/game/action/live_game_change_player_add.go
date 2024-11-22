package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerAddChange struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

func NewPlayerAddChange(userID uuid.UUID, playerID uuid.UUID, partySlot int) *PlayerAddChange {
	return &PlayerAddChange{UserID: userID, PlayerID: playerID, PartySlot: partySlot}
}

func (c PlayerAddChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerAdd
}
