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

func (s *PlayerTurnState) Update(game *game.LiveGameInstance, t time.Time) {

	var ready = game.GetActionLockedPlayerCount() == game.GetPlayerCount()
	var expired = s.Expired(t)

	if ready || expired {

		if expired {

			for _, party := range game.Parties {
				for _, player := range party.Players {
					if player.ActionsLocked == false {
						player.ActionLockIndex = party.GetActionLockedPlayerCount()
						player.ActionsLocked = true
						game.SendChange(NewPlayerLockActionChange(game.InstanceID, party.PartyIndex, player.PartySlot, player.ActionLockIndex))
					}
				}
			}

		}

		game.ActionChannel <- NewStateChangeAction(game.InstanceID, NewEnemyTurnState(game))

	}

}
