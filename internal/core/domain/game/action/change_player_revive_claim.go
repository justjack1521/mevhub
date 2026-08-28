package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

// PlayerReviveClaimChange announces a granted (or refreshed) revive claim so
// every client can mark the target as "being revived" and grey out its own
// revive option. Remaining is the claim window as granted.
type PlayerReviveClaimChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	// SourceID is who holds the claim; equal to PlayerID on a self-revive claim.
	SourceID   uuid.UUID
	PartyIndex int
	PartySlot  int
	Remaining  time.Duration
}

func NewPlayerReviveClaimChange(instanceID, playerID, sourceID uuid.UUID, partyIndex, partySlot int, remaining time.Duration) *PlayerReviveClaimChange {
	return &PlayerReviveClaimChange{InstanceID: instanceID, PlayerID: playerID, SourceID: sourceID, PartyIndex: partyIndex, PartySlot: partySlot, Remaining: remaining}
}

func (c PlayerReviveClaimChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerReviveClaim
}

// PlayerReviveClaimExpireChange announces that a claim decayed without a
// revive landing — the claimant crashed, or the item purchase failed. Clients
// clear the "being revived" indicator; the corpse is claimable again. A claim
// that ends in a revive emits PlayerReviveChange instead, never this.
type PlayerReviveClaimExpireChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	// SourceID is whose claim lapsed.
	SourceID   uuid.UUID
	PartyIndex int
	PartySlot  int
}

func NewPlayerReviveClaimExpireChange(instanceID, playerID, sourceID uuid.UUID, partyIndex, partySlot int) *PlayerReviveClaimExpireChange {
	return &PlayerReviveClaimExpireChange{InstanceID: instanceID, PlayerID: playerID, SourceID: sourceID, PartyIndex: partyIndex, PartySlot: partySlot}
}

func (c PlayerReviveClaimExpireChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierPlayerReviveClaimExpire
}
