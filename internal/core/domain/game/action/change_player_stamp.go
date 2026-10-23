package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerStampChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
	StampID    uuid.UUID
}

func NewPlayerStampChange(id uuid.UUID, party int, slot int, stamp uuid.UUID) *PlayerStampChange {
	return &PlayerStampChange{InstanceID: id, PartyIndex: party, PartySlot: slot, StampID: stamp}
}

func (c PlayerStampChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerStamp
}
