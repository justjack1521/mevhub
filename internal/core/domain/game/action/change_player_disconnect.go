package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerDisconnectChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerDisconnectChange(instanceID, playerID uuid.UUID, partyIndex, partySlot int) *PlayerDisconnectChange {
	return &PlayerDisconnectChange{InstanceID: instanceID, PlayerID: playerID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerDisconnectChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerDisconnect
}
