package game

// Change is the aggregate's outbound vocabulary. Not every change reaches a
// client — some are internal bookkeeping — so the wire ordinal is assigned
// where the bytes are published, not here.
type Change interface {
	Identifier() ChangeIdentifier
}

type ChangeIdentifier string

const (
	ChaneIdentifierPartyAdd          ChangeIdentifier = "party.add"
	ChangeIdentifierPlayerAdd        ChangeIdentifier = "player.add"
	ChangeIdentifierPlayerRemove     ChangeIdentifier = "player.remove"
	ChangeIdentifierPlayerReady      ChangeIdentifier = "player.ready"
	ChangeIdentifierPlayerDisconnect ChangeIdentifier = "player.disconnect"
	ChangeIdentifierPlayerReconnect  ChangeIdentifier = "player.reconnect"
	ChangeIdentifierPlayerDeath      ChangeIdentifier = "player.death"
	ChangeIdentifierPlayerRevive     ChangeIdentifier = "player.revive"
	ChangeIdentifierStateChange      ChangeIdentifier = "state.change"
	ChangeIdentifierEnqueueAction    ChangeIdentifier = "enqueue.action"
	ChangeIdentifierDequeueAction    ChangeIdentifier = "dequeue.action"
	ChangeIdentifierLockAction       ChangeIdentifier = "lock.action"
	ChangeIdentifierHPConsensus      ChangeIdentifier = "hp.consensus"
	ChangeIdentifierPlayerChat       ChangeIdentifier = "player.chat"
	ChangeIdentifierGameSync         ChangeIdentifier = "game.sync"
	// A granted revive claim and its unconsumed decay are both broadcast so
	// every client can show (and clear) a "being revived" indicator; a claim
	// that ends in a revive is signalled by the revive itself.
	ChangeIdentifierPlayerReviveClaim       ChangeIdentifier = "player.revive.claim"
	ChangeIdentifierPlayerReviveClaimExpire ChangeIdentifier = "player.revive.claim.expire"
)
