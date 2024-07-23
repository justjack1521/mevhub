package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type Instance struct {
	SysID        uuid.UUID
	PartyID      string
	Seed         int
	Options      *InstanceOptions
	State        InstanceState
	StartedAt    time.Time
	RegisteredAt time.Time
}

func NewGameInstance() *Instance {
	return &Instance{
		SysID:        uuid.NewV4(),
		PartyID:      fmt.Sprintf("%08d", rand.Intn(100000000)),
		Seed:         rand.Int(),
		State:        InstanceGamePendingState,
		RegisteredAt: time.Now().UTC(),
	}
}

type InstanceState int

const (
	InstanceGamePendingState InstanceState = 100
	InstanceGameStartedState InstanceState = 200
)
