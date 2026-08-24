package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

// restateGame broadcasts the whole of the current game state and records when it
// did so. It is the repair mechanism for a lost notification: rather than
// tracking what each client has seen and resending it, the game periodically
// says everything it knows, and any client that missed something converges on
// the next one.
//
// Every caller is already inside the single-writer loop, which is the only
// context where the aggregate may be read.
func restateGame(instance *game.LiveGameInstance) {
	instance.LastSyncAt = time.Now().UTC()
	instance.SendChange(NewGameStateSyncChange(instance))
}

// restateGamePeriodically emits a restatement if enough time has passed since
// the last one, from any source. Every state calls this from Update.
//
// A transition already restates, so in a game that keeps moving this rarely
// fires. It exists for the cases a transition cannot cover: the pending phase,
// which is entered by the factory rather than through transitionTo; a player
// reconnecting mid-turn, where nothing changes state; and a long player turn
// where a client could otherwise sit wrong until the turn ends.
func restateGamePeriodically(instance *game.LiveGameInstance, t time.Time) {
	if t.Sub(instance.LastSyncAt) < game.StateRestatePeriod {
		return
	}
	restateGame(instance)
}
