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
