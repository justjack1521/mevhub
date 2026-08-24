package action

import "mevhub/internal/core/domain/game"

// transitionTo swaps the game state, emits the StateChange notification, and
// restates the whole game behind it. Direct assignment is safe because states
// only run inside the game's single-writer loop, and it is synchronous — a
// boundary condition can never re-fire on a later tick while a transition is
// still queued.
//
// The restatement follows the StateChange rather than replacing it: the
// discrete notification is what a live client acts on immediately, and the
// restatement is what makes a client that missed it correct anyway. Emitting it
// second means it can only ever confirm what the StateChange already said.
func transitionTo(instance *game.LiveGameInstance, state game.State) {
	instance.State = state
	instance.SendChange(NewStateChange(instance.InstanceID, state))
	restateGame(instance)
}
