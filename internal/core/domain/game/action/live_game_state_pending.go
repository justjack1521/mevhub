package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

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

func NewPendingState(instance *game.LiveGameInstance) *PendingState {
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			player.Ready = false
		}
	}
	return &PendingState{
		StartTime:       time.Now().UTC(),
		MaxWaitDuration: game.PendingStateMaxWaitDuration,
	}
}

func (s *PendingState) Update(game *game.LiveGameInstance, t time.Time) {

	var expired = s.Expired(t)

	if expired {
		game.ActionChannel <- NewStateChangeAction(game.InstanceID, NewEndGameState(game))
		return
	}

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {
		game.ActionChannel <- NewStateChangeAction(game.InstanceID, NewPlayerTurnState(game))
	}
}
