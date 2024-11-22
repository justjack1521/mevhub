package game

import (
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

func NewPendingState(game *LiveGameInstance) *PendingState {
	for _, party := range game.Parties {
		for _, player := range party.Players {
			player.Ready = false
		}
	}
	return &PendingState{
		StartTime:       time.Now().UTC(),
		MaxWaitDuration: pendingStateMaxWaitDuration,
	}
}

func (s *PendingState) Update(game *LiveGameInstance, t time.Time) {

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
