package action

import (
	"errors"

	"mevhub/internal/core/domain/game"
)

// ErrNotPlayerTurn rejects a queue mutation outside the player turn. The
// per-player checks (dead, locked, queue cap) cannot express this on their own:
// in the pending phase every player starts alive and unlocked, so without this
// gate actions queued before everyone is ready are accepted and broadcast, then
// silently wiped when the first player turn resets the queues.
var ErrNotPlayerTurn = errors.New("actions can only be queued during the player turn")

// requirePlayerTurn is the state gate for enqueue, dequeue and lock. It lives
// here rather than on LivePlayer because the game package does not know the
// concrete states.
func requirePlayerTurn(instance *game.LiveGameInstance) error {
	if _, ok := instance.State.(*PlayerTurnState); ok {
		return nil
	}
	return ErrNotPlayerTurn
}
