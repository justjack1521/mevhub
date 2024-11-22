package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerDequeueActionChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerDequeueActionChange(instanceID uuid.UUID, partyIndex int, partySlot int) *PlayerDequeueActionChange {
	return &PlayerDequeueActionChange{InstanceID: instanceID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerDequeueActionChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierDequeueAction
}
