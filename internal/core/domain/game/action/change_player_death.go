package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerDeathChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerDeathChange(instanceID, playerID uuid.UUID, partyIndex, partySlot int) *PlayerDeathChange {
	return &PlayerDeathChange{InstanceID: instanceID, PlayerID: playerID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerDeathChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerDeath
}
