package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerDisconnectChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func NewPlayerDisconnectChange(instanceID uuid.UUID, playerID uuid.UUID) *PlayerDisconnectChange {
	return &PlayerDisconnectChange{InstanceID: instanceID, PlayerID: playerID}
}

func (c PlayerDisconnectChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerDisconnect
}
