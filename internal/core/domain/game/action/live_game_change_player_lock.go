package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerLockActionChange struct {
	InstanceID      uuid.UUID
	PartyIndex      int
	PartySlot       int
	ActionLockIndex int
}

func NewPlayerLockActionChange(id uuid.UUID, party int, slot int, index int) *PlayerLockActionChange {
	return &PlayerLockActionChange{InstanceID: id, PartyIndex: party, PartySlot: slot, ActionLockIndex: index}
}

func (c PlayerLockActionChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierLockAction
}
