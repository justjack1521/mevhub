package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

// pendingReadyTimeout bounds the initial load-in wait. Without it the only limit
// was PendingStateMaxWaitDuration (30 minutes), so one client whose ready never
// landed - dropped, or sent before the player was added - held the whole party
// for half an hour. Past it the game starts like the enemy turn's ready timeout
// does: nobody is removed, and a client still loading picks the turn up from the
// next sync frame (PLAYER_TURN clears its pre-turn gate). It counts from the
// game going live, before clients have even been told to load, so it has to
// cover a slow device's whole battle load.
const pendingReadyTimeout = time.Second * 60

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

	if instance.AllAlivePlayersReady() == false {
		if t.Sub(s.StartTime) <= pendingReadyTimeout {
			return
		}
		// Out of time, but hold for anyone in their reconnect grace, as the enemy
		// turn does at its boundary: the stall is bounded by the grace eviction,
		// and starting would hand a returning player a turn already running.
		if instance.HasDisconnectedPlayers() {
			return
		}
	}

	instance.StartedAt = t
	transitionTo(instance, NewPlayerTurnState(instance))
}
