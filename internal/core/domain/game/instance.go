package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Instance struct {
	SysID        uuid.UUID
	Seed         int
	LobbyIDs     []uuid.UUID
	Options      *InstanceOptions
	State        InstanceState
	RegisteredAt time.Time
}

type InstanceState int

const (
	InstanceGamePendingState InstanceState = 100
	InstanceGameStartedState InstanceState = 200
)

type Summary struct {
	SysID   uuid.UUID
	Seed    int
	Parties []PartySummary
}
