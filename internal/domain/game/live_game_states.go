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

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {

		for _, player := range game.Players {
			player.ActionsLocked = false
		}

		game.ActionChannel <- &StateChangeAction{
			InstanceID: game.InstanceID,
			State:      game.NewPlayerTurnState(),
		}

	}
}

type PlayerTurnState struct {
	StartTime    time.Time
	TurnDuration time.Duration
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
			State:      game.NewEnemyTurnState(),
		}
	}

}

type EnemyTurnState struct {
	QueuedActions []*PlayerActionQueue
}

func (s *EnemyTurnState) Update(game *LiveGameInstance, t time.Time) {

}
