package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerReconnectChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerReconnectChange(instanceID, playerID uuid.UUID, partyIndex, partySlot int) *PlayerReconnectChange {
	return &PlayerReconnectChange{InstanceID: instanceID, PlayerID: playerID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerReconnectChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerReconnect
}
