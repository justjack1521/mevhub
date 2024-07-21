package game

import (
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type LiveGameInstance struct {
	ActionChannel chan Action
	ChangeChannel chan Change
	Players       map[uuid.UUID]*LivePlayer
	Mu            sync.RWMutex
}

func NewLiveGameInstance() *LiveGameInstance {
	return &LiveGameInstance{
		ActionChannel: make(chan Action),
		ChangeChannel: make(chan Change),
		Players:       make(map[uuid.UUID]*LivePlayer),
	}
}

func (game *LiveGameInstance) WatchActions() {
	for {
		action := <-game.ActionChannel
		game.Mu.Lock()
		action.Perform(game)
		game.Mu.Unlock()
	}
}

func (game *LiveGameInstance) SendChange(change Change) {
	select {
	case game.ChangeChannel <- change:
	default:
	}
}

type LivePlayer struct {
	UserID     uuid.UUID
	PlayerID   uuid.UUID
	PartySlot  int
	Ready      bool
	LastAction time.Time
}
