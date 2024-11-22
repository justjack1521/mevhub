package game

import uuid "github.com/satori/go.uuid"

type StateChange struct {
	InstanceID uuid.UUID
	State      State
}

func NewStateChange(id uuid.UUID, state State) *StateChange {
	return &StateChange{InstanceID: id, State: state}
}

func (c StateChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierStateChange
}
