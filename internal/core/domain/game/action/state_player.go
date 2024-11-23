package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

type PlayerTurnState struct {
	StartTime    time.Time
	TurnDuration time.Duration
}

func NewPlayerTurnState(game *game.LiveGameInstance) *PlayerTurnState {
	for _, party := range game.Parties {
		for _, player := range party.Players {
			player.ActionsLocked = false
		}
	}
	return &PlayerTurnState{
		StartTime:    time.Now().UTC(),
		TurnDuration: game.PartyOptions.PlayerTurnDuration,
	}
}

func (s *PlayerTurnState) Expired(t time.Time) bool {
	if s.TurnDuration == 0 {
		return false
	}
	var difference = t.Sub(s.StartTime)
	return difference > s.TurnDuration
}

func (s *PlayerTurnState) Update(instance *game.LiveGameInstance, t time.Time) {

	var ready = instance.GetActionLockedPlayerCount() == instance.GetPlayerCount()
	var expired = s.Expired(t)

	if ready || expired {

		if expired {

			for _, party := range instance.Parties {
				for _, player := range party.Players {
					if player.ActionsLocked == false {
						player.ActionLockIndex = party.GetActionLockedPlayerCount()
						player.ActionsLocked = true
						instance.SendChange(NewPlayerLockActionChange(instance.InstanceID, party.PartyIndex, player.PartySlot, player.ActionLockIndex))
					}
				}
			}

		}

		instance.ActionChannel <- NewStateChangeAction(instance.InstanceID, NewEnemyTurnState(instance))

	}

}
