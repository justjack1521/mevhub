package game

import (
	"time"
)

type State interface {
	Update(game *LiveGameInstance, t time.Time)
}

type PendingState struct {
	StartTime time.Time
}

func (s *PendingState) Update(game *LiveGameInstance, t time.Time) {
	if game.GetReadyPlayerCount() == game.GetPlayerCount() {

		for _, player := range game.Players {
			player.ActionsLocked = false
		}

		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State: &PlayerTurnState{
				StartTime: time.Now().UTC(),
			},
		}

	}
}

type PlayerTurnState struct {
	StartTime time.Time
}

func (s *PlayerTurnState) Update(game *LiveGameInstance, t time.Time) {

	var difference = t.Sub(s.StartTime)

	var ready = game.GetActionLockedPlayerCount() == game.GetPlayerCount()
	var expired = difference > game.PlayerTurnDuration

	if ready || expired {

		for _, player := range game.Players {
			player.ActionsLocked = true
			var change = PlayerLockActionChange{
				InstanceID: game.InstanceID,
				PartySlot:  player.PartySlot,
			}
			game.SendChange(change)
		}

		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      &EnemyTurnState{},
		}
	}

}

type EnemyTurnState struct {
}

func (s *EnemyTurnState) Update(game *LiveGameInstance, t time.Time) {

}
