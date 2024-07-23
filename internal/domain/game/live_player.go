package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type LivePlayer struct {
	UserID          uuid.UUID
	PlayerID        uuid.UUID
	PartySlot       int
	Ready           bool
	ActionsLocked   bool
	ActionLockIndex int
	MaxActionCount  int
	Actions         []*PlayerAction
	LastAction      time.Time
}

func (p *LivePlayer) CanEnqueueAction() bool {
	return p.ActionsLocked == false && len(p.Actions) < p.MaxActionCount
}

func (p *LivePlayer) CanDequeueAction() bool {
	return p.ActionsLocked == false && len(p.Actions) > 0
}

func (p *LivePlayer) DequeueAction() bool {
	if p.CanDequeueAction() == false {
		return false
	}
	p.Actions = p.Actions[:len(p.Actions)-1]
	return true
}

func (p *LivePlayer) EnqueueAction(action *PlayerAction) bool {

	if p.CanEnqueueAction() == false {
		return false
	}

	p.Actions = append(p.Actions, action)
	return true

}

type PlayerAction struct {
	Target     int
	ActionType PlayerActionType
	SlotIndex  int
	ElementID  uuid.UUID
}

type PlayerActionType int

const (
	PlayerActionTypeNone = iota
	PlayerActionTypeNormalAttack
	PlayerActionTypeAbilityCast
	PlayerActionTypeElementDrive
)
