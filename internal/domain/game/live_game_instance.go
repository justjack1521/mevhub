package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	GameInstanceTickPeriod = time.Millisecond * 20
)

type LiveGameInstance struct {
	InstanceID         uuid.UUID
	ActionChannel      chan Action
	ChangeChannel      chan Change
	Players            map[uuid.UUID]*LivePlayer
	State              State
	PlayerTurnDuration time.Duration
	GameDuration       time.Duration
	MaxPlayerCount     int
}

func NewLiveGameInstance() *LiveGameInstance {
	return &LiveGameInstance{
		ActionChannel: make(chan Action),
		ChangeChannel: make(chan Change),
		Players:       make(map[uuid.UUID]*LivePlayer),
		State:         &PendingState{StartTime: time.Now().UTC()},
	}
}

func (game *LiveGameInstance) GetPlayerCount() int {
	return len(game.Players)
}

func (game *LiveGameInstance) GetReadyPlayerCount() int {
	var count int
	for _, player := range game.Players {
		if player.Ready {
			count++
		}
	}
	return count
}

func (game *LiveGameInstance) GetActionLockedPlayerCount() int {
	var count int
	for _, player := range game.Players {
		if player.ActionsLocked {
			count++
		}
	}
	return count
}

func (game *LiveGameInstance) Tick() {

	ticker := time.NewTicker(GameInstanceTickPeriod)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			game.State.Update(game, t)
		}
	}

}

func (game *LiveGameInstance) WatchActions() {
	for {
		action := <-game.ActionChannel
		action.Perform(game)
	}
}

func (game *LiveGameInstance) SendChange(change Change) {
	select {
	case game.ChangeChannel <- change:
	default:
	}
}
