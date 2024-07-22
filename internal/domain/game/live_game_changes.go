package game

import (
	uuid "github.com/satori/go.uuid"
)

type Change interface {
}

type PlayerAddChange struct {
	Change
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

type PlayerReadyChange struct {
	Change
	InstanceID uuid.UUID
	PartySlot  int
}

type GameStateChange struct {
	Change
	InstanceID uuid.UUID
}

type PlayerEnqueueActionChange struct {
	Change
	InstanceID uuid.UUID
	PartySlot  int
	ActionType PlayerActionType
	SlotIndex  int
	Target     int
	ElementID  uuid.UUID
}

type PlayerDequeueActionChange struct {
	Change
	InstanceID uuid.UUID
	PartySlot  int
}

type PlayerLockActionChange struct {
	Change
	InstanceID uuid.UUID
	PartySlot  int
}
