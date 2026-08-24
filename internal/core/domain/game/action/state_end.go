package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

type EndGameState struct {
}

// Update is never called: NewEndGameState calls instance.End(), and the game
// loop skips State.Update once the game has ended. It is deliberately left
// empty rather than given a periodic restatement — there would be nothing to
// drive it. The transition into this state restates once on the way in, which
// is the only chance a client gets to learn the game is over.
func (s *EndGameState) Update(instance *game.LiveGameInstance, t time.Time) {

}

func NewEndGameState(instance *game.LiveGameInstance) *EndGameState {

	instance.End()

	return &EndGameState{}
}
