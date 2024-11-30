package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerRemoveChange struct {
	InstanceID uuid.UUID
	UserID     uuid.UUID
	PlayerID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerRemoveChange(instance, user, player uuid.UUID, party int, slot int) *PlayerRemoveChange {
	return &PlayerRemoveChange{InstanceID: instance, UserID: user, PlayerID: player, PartyIndex: party, PartySlot: slot}
}

func (c PlayerRemoveChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerRemove
}
