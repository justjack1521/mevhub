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
	ChangeIdentifierStateChange      ChangeIdentifier = "state.change"
	ChangeIdentifierEnqueueAction    ChangeIdentifier = "enqueue.action"
	ChangeIdentifierDequeueAction    ChangeIdentifier = "dequeue.action"
	ChangeIdentifierLockAction       ChangeIdentifier = "lock.action"
	ChangeIdentifierHPConsensus      ChangeIdentifier = "hp.consensus"
	ChangeIdentifierGameSync         ChangeIdentifier = "game.sync"
)
