package action

import "mevhub/internal/core/domain/game"

// transitionTo swaps the game state and emits the StateChange notification.
// Direct assignment is safe because states only run inside the game's
// single-writer loop, and it is synchronous — a boundary condition can never
// re-fire on a later tick while a transition is still queued.
func transitionTo(instance *game.LiveGameInstance, state game.State) {
	instance.State = state
	instance.SendChange(NewStateChange(instance.InstanceID, state))
}
