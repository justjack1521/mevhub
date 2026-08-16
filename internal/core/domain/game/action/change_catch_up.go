package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

// CatchUpResult reports what a replay actually covered. When there was nothing
// to send the range is inverted (FromSequence > ToSequence) rather than zeroed,
// so "you are up to date" is distinguishable from "here is the game from the
// top". A FromSequence higher than the client asked for means the rest had
// already been evicted.
type CatchUpResult struct {
	FromSequence    uint64
	ToSequence      uint64
	TurnRemainingMs int64
}

// CatchUpChange asks the publisher to replay its notification log to a single
// player. Unlike every other change it carries no game state of its own: the
// payloads it stands for were marshalled once when they were first published,
// and the log holds them.
//
// It travels the change channel rather than being served directly from the gRPC
// goroutine so the replayed bytes leave on the same goroutine as live traffic.
// Anything the loop emits while the replay is in flight is queued behind it by
// construction, and so can never overtake it.
type CatchUpChange struct {
	InstanceID     uuid.UUID
	TargetPlayerID uuid.UUID
	// FromSequence is the first ordinal the client wants, not the last it has.
	FromSequence uint64
	// TurnRemainingMs is the one thing replay cannot reconstruct. Discrete
	// events replay perfectly; elapsed wall-clock does not, so the remaining
	// turn time is sampled here, inside the loop, and rides the response. -1
	// means the turn has no timer.
	TurnRemainingMs int64
	result          chan<- CatchUpResult
}

func NewCatchUpChange(instance uuid.UUID, target uuid.UUID, from uint64, remaining int64, result chan<- CatchUpResult) *CatchUpChange {
	return &CatchUpChange{
		InstanceID:      instance,
		TargetPlayerID:  target,
		FromSequence:    from,
		TurnRemainingMs: remaining,
		result:          result,
	}
}

// Report hands the covered range back to the caller waiting on the RPC. The
// channel is buffered and the send is non-blocking, so a caller that has
// already given up (its context expired) cannot wedge the publisher.
func (c *CatchUpChange) Report(result CatchUpResult) {
	if c.result == nil {
		return
	}
	select {
	case c.result <- result:
	default:
	}
}

func (c CatchUpChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierCatchUp
}

// turnRemainingMs samples how much of the current player turn is left. It must
// be called from inside the single-writer loop, which is the only context where
// the state may be read. Anything that is not a timed player turn reports -1,
// the established "no timer" convention — 0 would read as "already over".
func turnRemainingMs(instance *game.LiveGameInstance) int64 {
	state, ok := instance.State.(*PlayerTurnState)
	if !ok || state.TurnDuration <= 0 {
		return -1
	}
	var remaining = state.TurnDuration - time.Since(state.StartTime)
	if remaining < 0 {
		remaining = 0
	}
	return remaining.Milliseconds()
}
