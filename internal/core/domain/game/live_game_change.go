package game

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
