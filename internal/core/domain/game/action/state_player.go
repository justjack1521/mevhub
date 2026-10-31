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
			// A fresh turn starts from a clean slate: without these resets the
			// action queue accumulated across turns and a stale lock index
			// from a larger previous turn could overrun the next queue.
			player.ActionsLocked = false
			player.ActionLockIndex = 0
			player.Actions = nil
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
	return difference >= s.TurnDuration
}

func (s *PlayerTurnState) Update(instance *game.LiveGameInstance, t time.Time) {

	evictExpiredDisconnectedPlayers(instance, t)
	evictExpiredDeadPlayers(instance, t)
	expireLapsedReviveClaims(instance, t)
	restateGamePeriodically(instance, t)

	// An abandoned game must end, not spin: with zero players the ready count
	// comparison below would be trivially true and bounce this game between
	// turn states forever.
	if instance.GetPlayerCount() == 0 {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	// Out of time ends the game anywhere in the player turn. The action queue
	// is only dispatched on the way into the enemy turn, so no client has
	// played anything out yet and they all leave from the same point. Once it
	// is dispatched the end waits for the enemy -> player boundary instead.
	if instance.TimeExpired(t) {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	// Everyone dead is a wipe, and this turn is the party's chance to get back
	// up: anyone may revive themselves or another. Nobody alive means nobody to
	// lock, and the lock check below would pass vacuously and run an enemy turn
	// over corpses, so hold here. The turn timer running out with the party
	// still down ends the game; so does the dead-player kick, or everyone
	// quitting, when the quest has no turn timer.
	if instance.GetAlivePlayerCount() == 0 {
		if s.Expired(t) {
			transitionTo(instance, NewEndGameState(instance))
		}
		return
	}

	var ready = instance.AllAlivePlayersLocked()
	var expired = s.Expired(t)

	if ready || expired {

		// We are at the player -> enemy boundary. Hold the transition while any
		// player is within their reconnect grace period so the turn does not
		// end without them; the turn itself continues to accept actions. The
		// stall lifts on reconnect or timeout removal.
		if instance.HasDisconnectedPlayers() {
			return
		}

		if expired {

			for _, party := range instance.Parties {
				for _, player := range party.Players {
					// A dead player is exempt from locking, not force-locked:
					// they cannot act, and the enemy turn skips them anyway.
					if player.Dead {
						continue
					}
					if player.ActionsLocked == false {
						player.ActionLockIndex = party.NextActionLockIndex()
						player.ActionsLocked = true
						instance.SendChange(NewPlayerLockActionChange(instance.InstanceID, party.PartyIndex, player.PartySlot, player.ActionLockIndex))
					}
				}
			}

		}

		transitionTo(instance, NewEnemyTurnState(instance))

	}

}
