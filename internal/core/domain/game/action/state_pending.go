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

func (s *PendingState) Update(instance *game.LiveGameInstance, t time.Time) {

	evictExpiredDisconnectedPlayers(instance, t)
	evictExpiredDeadPlayers(instance, t)
	expireLapsedReviveClaims(instance, t)
	restateGamePeriodically(instance, t)

	var expired = s.Expired(t)

	if expired {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	if instance.GetPlayerCount() == 0 {
		return
	}

	// Alive-scoped like the turn boundaries: a death should not be reportable
	// before the battle starts, but if one ever is, the game must neither wait
	// on the corpse nor start with nobody able to act.
	if instance.GetAlivePlayerCount() == 0 {
		return
	}

	if instance.AllAlivePlayersReady() {
		transitionTo(instance, NewPlayerTurnState(instance))
	}
}
