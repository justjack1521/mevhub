package game

import (
	"errors"
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

var (
	ErrPlayerActionsLocked = errors.New("player actions locked")
	ErrPlayerActionsFull   = errors.New("player actions full")
	ErrPlayerActionsEmpty  = errors.New("player actions empty")
)

func (p *LivePlayer) CanEnqueueAction() error {

	if p.ActionsLocked {
		return ErrPlayerActionsLocked
	}

	if len(p.Actions) >= p.MaxActionCount && p.MaxActionCount > 0 {
		return ErrPlayerActionsFull
	}

	return nil
}

func (p *LivePlayer) CanDequeueAction() error {
	if p.ActionsLocked {
		return ErrPlayerActionsLocked
	}

	if len(p.Actions) == 0 {
		return ErrPlayerActionsEmpty
	}

	return nil
}

func (p *LivePlayer) DequeueAction() error {
	if err := p.CanDequeueAction(); err != nil {
		return err
	}
	p.Actions = p.Actions[:len(p.Actions)-1]
	return nil
}

func (p *LivePlayer) EnqueueAction(action *PlayerAction) error {

	if err := p.CanEnqueueAction(); err != nil {
		return err
	}

	p.Actions = append(p.Actions, action)
	return nil

}

type PlayerActionQueue struct {
	PlayerID uuid.UUID
	Actions  []*PlayerAction
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
