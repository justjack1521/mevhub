package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

var (
	ErrFailedCatchUpPlayer = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to catch up player %s: %w", player, err)
	}
)

// CatchUpAction asks the game to replay the notifications a player missed. It
// mutates nothing; it exists purely so the request enters the pipeline at the
// same point as every other client call and leaves in order behind whatever the
// loop has already emitted.
type CatchUpAction struct {
	GameID   uuid.UUID
	PlayerID uuid.UUID
	// LastSequence is the highest ordinal the client has applied. Zero means it
	// has nothing and wants the game from the top.
	LastSequence uint64
	Result       chan<- CatchUpResult
}

func NewCatchUpAction(gameID uuid.UUID, playerID uuid.UUID, last uint64, result chan<- CatchUpResult) *CatchUpAction {
	return &CatchUpAction{GameID: gameID, PlayerID: playerID, LastSequence: last, Result: result}
}

func (a *CatchUpAction) Perform(instance *game.LiveGameInstance) error {

	if instance.PlayerExists(a.PlayerID) == false {
		// Answer anyway: the caller is blocked on the reply, and an empty range
		// releases it now rather than at its context deadline.
		a.report(CatchUpResult{FromSequence: a.LastSequence + 1, ToSequence: a.LastSequence, TurnRemainingMs: -1})
		return ErrFailedCatchUpPlayer(a.PlayerID, game.ErrPlayerNotInGame)
	}

	instance.SendChange(NewCatchUpChange(instance.InstanceID, a.PlayerID, a.LastSequence+1, turnRemainingMs(instance), a.Result))
	return nil

}

func (a *CatchUpAction) report(result CatchUpResult) {
	if a.Result == nil {
		return
	}
	select {
	case a.Result <- result:
	default:
	}
}
