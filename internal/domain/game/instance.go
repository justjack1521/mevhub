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
	Seed         int64
	Options      InstanceOptions
	State        InstanceState
	RegisteredAt time.Time
}

func NewGameInstance() *Instance {
	return &Instance{
		SysID:        uuid.NewV4(),
		PartyID:      fmt.Sprintf("%08d", rand.Intn(100000000)),
		Seed:         rand.Int63(),
		State:        InstanceGamePendingState,
		RegisteredAt: time.Now().UTC(),
	}
}

type InstanceState int

const (
	InstanceGamePendingState InstanceState = 200
)
