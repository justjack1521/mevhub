package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type StateChange struct {
	InstanceID uuid.UUID
	State      game.State
}

func NewStateChange(id uuid.UUID, state game.State) *StateChange {
	return &StateChange{InstanceID: id, State: state}
}

func (c StateChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierStateChange
}
