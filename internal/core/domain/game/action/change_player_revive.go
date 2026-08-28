package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerReviveChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	// SourceID is who performed the revive; equal to PlayerID on a self-revive.
	SourceID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerReviveChange(instanceID, playerID, sourceID uuid.UUID, partyIndex, partySlot int) *PlayerReviveChange {
	return &PlayerReviveChange{InstanceID: instanceID, PlayerID: playerID, SourceID: sourceID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerReviveChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerRevive
}
