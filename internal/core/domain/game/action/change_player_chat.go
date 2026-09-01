package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerChatChange struct {
	InstanceID uuid.UUID
	PartyIndex int
	PartySlot  int
	Message    string
}

func NewPlayerChatChange(id uuid.UUID, party int, slot int, message string) *PlayerChatChange {
	return &PlayerChatChange{InstanceID: id, PartyIndex: party, PartySlot: slot, Message: message}
}

func (c PlayerChatChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerChat
}
