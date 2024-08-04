package game

import (
	"time"
)

const (
	pendingStateMaxWaitDuration = time.Minute * 1
)

type State interface {
	Update(game *LiveGameInstance, t time.Time)
}

type PendingState struct {
	StartTime       time.Time
	MaxWaitDuration time.Duration
}

func (s *PendingState) Expired(t time.Time) bool {
	if s.MaxWaitDuration == 0 {
		return false
	}
	var difference = t.Sub(s.StartTime)
	return difference > s.MaxWaitDuration
}

func NewPendingState(game *LiveGameInstance) *PendingState {
	for _, player := range game.Players {
		player.Ready = false
	}
	return &PendingState{
		StartTime:       time.Now().UTC(),
		MaxWaitDuration: pendingStateMaxWaitDuration,
	}
}

func (s *PendingState) Update(game *LiveGameInstance, t time.Time) {

	var expired = s.Expired(t)

	if expired {

		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      NewEndGameState(game),
		}

		return
	}

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {
		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      NewPlayerTurnState(game),
		}
	}

}

type EndGameState struct {
}

func (s *EndGameState) Update(game *LiveGameInstance, t time.Time) {

}

func NewEndGameState(game *LiveGameInstance) *EndGameState {

	game.Ended = true
	game.EndedAt = time.Now().UTC()

	return &EndGameState{}
}

type PlayerTurnState struct {
	StartTime    time.Time
	TurnDuration time.Duration
}

func NewPlayerTurnState(game *LiveGameInstance) *PlayerTurnState {
	for _, player := range game.Players {
		player.ActionsLocked = false
	}
	return &PlayerTurnState{
		StartTime:    time.Now().UTC(),
		TurnDuration: game.PlayerTurnDuration,
	}
}

func (s *PlayerTurnState) Expired(t time.Time) bool {
	if s.TurnDuration == 0 {
		return false
	}
	var difference = t.Sub(s.StartTime)
	return difference > s.TurnDuration
}

func (s *PlayerTurnState) Update(game *LiveGameInstance, t time.Time) {

	var ready = game.GetActionLockedPlayerCount() == game.GetPlayerCount()
	var expired = s.Expired(t)

	if ready || expired {

		if expired {
			for _, player := range game.Players {
				if player.ActionsLocked == false {
					player.ActionLockIndex = game.GetActionLockedPlayerCount()
					player.ActionsLocked = true

					var change = PlayerLockActionChange{
						InstanceID:      game.InstanceID,
						PartySlot:       player.PartySlot,
						ActionLockIndex: player.ActionLockIndex,
					}

					game.SendChange(change)
				}
			}
		}

		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      NewEnemyTurnState(game),
		}
	}

}

type EnemyTurnState struct {
	QueuedActions []*PlayerActionQueue
}

func NewEnemyTurnState(game *LiveGameInstance) *EnemyTurnState {
	var state = &EnemyTurnState{
		QueuedActions: make([]*PlayerActionQueue, len(game.Players)),
	}
	for _, player := range game.Players {
		player.Ready = false
		state.QueuedActions[player.ActionLockIndex] = &PlayerActionQueue{
			PlayerID: player.PlayerID,
			Actions:  player.Actions,
		}
	}
	return state
}

func (s *EnemyTurnState) Update(game *LiveGameInstance, t time.Time) {

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {
		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      NewPlayerTurnState(game),
		}
	}

}
