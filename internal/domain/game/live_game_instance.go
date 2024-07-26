package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	StateTickPeriod = time.Millisecond * 250
)

type LiveGameInstance struct {
	InstanceID         uuid.UUID
	ActionChannel      chan Action
	ChangeChannel      chan Change
	ErrorChannel       chan error
	Players            map[uuid.UUID]*LivePlayer
	State              State
	PlayerTurnDuration time.Duration
	GameDuration       time.Duration
	MaxPlayerCount     int
}

func NewLiveGameInstance(source *Instance) *LiveGameInstance {
	var game = &LiveGameInstance{
		InstanceID:         source.SysID,
		ActionChannel:      make(chan Action),
		ChangeChannel:      make(chan Change),
		ErrorChannel:       make(chan error),
		Players:            make(map[uuid.UUID]*LivePlayer),
		PlayerTurnDuration: source.Options.PlayerTurnDuration,
		GameDuration:       source.Options.MaxRunTime,
		MaxPlayerCount:     source.Options.MaxPlayerCount,
	}
	game.State = game.NewPendingState()
	return game
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

	ticker := time.NewTicker(StateTickPeriod)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			game.State.Update(game, t)
		}
	}

}

var (
	ErrFailedPerformAction = func(id uuid.UUID, err error) error {
		return fmt.Errorf("live game %s failed to perform action: %w", id, err)
	}
)

func (game *LiveGameInstance) WatchActions() {
	for {
		action := <-game.ActionChannel
		if err := action.Perform(game); err != nil {
			game.ErrorChannel <- ErrFailedPerformAction(game.InstanceID, err)
		}
	}
}

func (game *LiveGameInstance) SendChange(change Change) {
	game.ChangeChannel <- change
}

func (game *LiveGameInstance) NewPendingState() *PendingState {
	return &PendingState{StartTime: time.Now().UTC()}
}

func (game *LiveGameInstance) NewPlayerTurnState() *PlayerTurnState {
	return &PlayerTurnState{
		StartTime:    time.Now().UTC(),
		TurnDuration: game.PlayerTurnDuration,
	}
}

func (game *LiveGameInstance) NewEnemyTurnState() *EnemyTurnState {
	var state = &EnemyTurnState{QueuedActions: make([]*PlayerActionQueue, len(game.Players))}
	for _, player := range game.Players {
		state.QueuedActions[player.ActionLockIndex] = &PlayerActionQueue{
			PlayerID: player.PlayerID,
			Actions:  player.Actions,
		}
	}
	return state
}

func (game *LiveGameInstance) GetPlayer(id uuid.UUID) (*LivePlayer, error) {
	player, exists := game.Players[id]
	if exists == false {
		return nil, ErrPlayerNotInGame
	}
	return player, nil
}
