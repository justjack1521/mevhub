package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	StateTickPeriod = time.Millisecond * 250
)

var (
	ErrFailedPerformAction = func(id uuid.UUID, err error) error {
		return fmt.Errorf("live game %s failed to perform action: %w", id, err)
	}
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
	Ended              bool
	EndedAt            time.Time
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
	game.State = NewPendingState(game)
	return game
}

func (game *LiveGameInstance) GetPlayerCount() int {
	return len(game.Players)
}

func (game *LiveGameInstance) PlayerExists(id uuid.UUID) bool {
	_, exists := game.Players[id]
	return exists
}

func (game *LiveGameInstance) GetPlayer(id uuid.UUID) (*LivePlayer, error) {
	player, exists := game.Players[id]
	if exists == false {
		return nil, ErrPlayerNotInGame
	}
	return player, nil
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

func (game *LiveGameInstance) RemovePlayer(id uuid.UUID) error {
	if game.PlayerExists(id) == false {
		return ErrPlayerNotInGame
	}
	delete(game.Players, id)
	return nil
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
