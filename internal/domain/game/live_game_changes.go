package game

import (
	uuid "github.com/satori/go.uuid"
)

type Change interface {
	Identifier() ChangeIdentifier
}

type ChangeIdentifier string

const (
	ChangeIdentifierPlayerAdd        ChangeIdentifier = "player.add"
	ChangeIdentifierPlayerReady      ChangeIdentifier = "player.ready"
	ChangeIdentifierPlayerDisconnect ChangeIdentifier = "player.disconnect"
	ChangeIdentifierStateChange      ChangeIdentifier = "state.change"
	ChangeIdentifierEnqueueAction    ChangeIdentifier = "enqueue.action"
	ChangeIdentifierDequeueAction    ChangeIdentifier = "dequeue.action"
	ChangeIdentifierLockAction       ChangeIdentifier = "lock.action"
)

type PlayerAddChange struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

func (c PlayerAddChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierPlayerAdd
}

type PlayerReadyChange struct {
	InstanceID uuid.UUID
	PartySlot  int
}

func (c PlayerReadyChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierPlayerReady
}

type PlayerDisconnectChange struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (c PlayerDisconnectChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierPlayerDisconnect
}

type StateChange struct {
	InstanceID uuid.UUID
	State      State
}

func (c StateChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierStateChange
}

type PlayerEnqueueActionChange struct {
	InstanceID uuid.UUID
	PartySlot  int
	ActionType PlayerActionType
	SlotIndex  int
	Target     int
	ElementID  uuid.UUID
}

func (c PlayerEnqueueActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierEnqueueAction
}

type PlayerDequeueActionChange struct {
	InstanceID uuid.UUID
	PartySlot  int
}

func (c PlayerDequeueActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierDequeueAction
}

type PlayerLockActionChange struct {
	InstanceID      uuid.UUID
	PartySlot       int
	ActionLockIndex int
}

func (c PlayerLockActionChange) Identifier() ChangeIdentifier {
	return ChangeIdentifierLockAction
}
