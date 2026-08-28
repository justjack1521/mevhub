package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

// ReviveClaimDenyReason enumerates why a revive claim was refused. The values
// are kept in sync with the protomulti deny-reason enum so the adapter can
// cast directly.
type ReviveClaimDenyReason int

const (
	ReviveClaimDenyReasonNone ReviveClaimDenyReason = iota
	ReviveClaimDenyReasonTargetNotFound
	ReviveClaimDenyReasonTargetAlive
	ReviveClaimDenyReasonAlreadyClaimed
	// ReviveClaimDenyReasonUnavailable is never produced here: the application
	// layer stamps it when no reply arrives at all (game not hosted, shut down
	// mid-route, or the loop too far behind), so the claimant fails closed and
	// spends nothing.
	ReviveClaimDenyReasonUnavailable
)

// PlayerReviveClaimResult is the synchronous answer to a claim. Remaining is
// the window that matters to the caller: on a grant, how long they have to
// consume the item and call revive before the claim decays; on an
// already-claimed denial, how long until the standing claim decays and the
// corpse can be claimed again.
type PlayerReviveClaimResult struct {
	Granted   bool
	Reason    ReviveClaimDenyReason
	Remaining time.Duration
}

// PlayerReviveClaimAction takes the decaying exclusive right to revive a dead
// player. Unlike every other action it answers its caller: reviving costs
// items consumed through a separate service, so the claimant must know the
// claim held BEFORE spending — a fire-and-forget claim would protect nothing.
// The result is sent on Response, which must be buffered; the send never
// blocks the game loop, and a caller that stopped listening just misses an
// answer it no longer wants.
//
// A denial is an answer, not a failure, so Perform never returns an error —
// losing the race to another claimant puts nothing on the error channel.
type PlayerReviveClaimAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	SourceID   uuid.UUID
	ClaimTime  time.Time
	Response   chan PlayerReviveClaimResult
}

func NewPlayerReviveClaimAction(instanceID uuid.UUID, playerID uuid.UUID, sourceID uuid.UUID, claimTime time.Time, response chan PlayerReviveClaimResult) *PlayerReviveClaimAction {
	return &PlayerReviveClaimAction{InstanceID: instanceID, PlayerID: playerID, SourceID: sourceID, ClaimTime: claimTime, Response: response}
}

func (a *PlayerReviveClaimAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		a.reply(PlayerReviveClaimResult{Reason: ReviveClaimDenyReasonTargetNotFound})
		return nil
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		a.reply(PlayerReviveClaimResult{Reason: ReviveClaimDenyReasonTargetNotFound})
		return nil
	}

	// Claiming the living is the stale-client race this exists to stop: the
	// target was revived moments ago and this claimant has not seen it yet.
	// Denying here is what saves their item.
	if player.Dead == false {
		a.reply(PlayerReviveClaimResult{Reason: ReviveClaimDenyReasonTargetAlive})
		return nil
	}

	if player.HasActiveReviveClaim(a.ClaimTime) && uuid.Equal(player.ReviveClaimSource, a.SourceID) == false {
		a.reply(PlayerReviveClaimResult{
			Reason:    ReviveClaimDenyReasonAlreadyClaimed,
			Remaining: player.ReviveClaimExpiry.Sub(a.ClaimTime),
		})
		return nil
	}

	// Free, expired, or re-claimed by the same claimant (an idempotent retry
	// refreshes the window rather than deadlocking against itself).
	player.ReviveClaimSource = a.SourceID
	player.ReviveClaimExpiry = a.ClaimTime.Add(game.ReviveClaimDuration)

	// Broadcast the grant (a refresh re-broadcasts with the new window) so
	// every client can show "being revived" and stop offering the revive.
	instance.SendChange(NewPlayerReviveClaimChange(instance.InstanceID, a.PlayerID, a.SourceID, party.PartyIndex, player.PartySlot, game.ReviveClaimDuration))

	a.reply(PlayerReviveClaimResult{Granted: true, Remaining: game.ReviveClaimDuration})
	return nil
}

// reply never blocks the game loop: Response is buffered by the caller, and if
// the caller has already timed out and gone, the answer is dropped.
func (a *PlayerReviveClaimAction) reply(result PlayerReviveClaimResult) {
	select {
	case a.Response <- result:
	default:
	}
}
