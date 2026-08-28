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
	// MaxActionCount caps the queue in CanEnqueueAction; 0 means unlimited.
	// Nothing populates it from loadout data yet, so every live player is
	// currently uncapped.
	MaxActionCount int
	Actions        []*PlayerAction
	LastAction     time.Time
	Disconnected   bool
	DisconnectTime time.Time
	// Dead marks a player who has fallen in battle. A dead player cannot act
	// (enqueue/dequeue) but stays in the game and may be revived — by another
	// player or by themselves. DeathTime anchors the kick sweep: a player dead
	// longer than the game's DeadPlayerKickDuration with no revive is removed.
	Dead      bool
	DeathTime time.Time
	// ReviveClaimSource and ReviveClaimExpiry form the decaying revive claim.
	// Reviving costs items consumed through a separate service (BattleRevive on
	// MeviusGameService), so before consuming anything a client must claim the
	// exclusive right to revive this corpse; a second claimant is denied and
	// spends nothing. The claim decays at ReviveClaimExpiry so an abandoned
	// claim (claimant crashed, item purchase failed) never leaves a corpse
	// unrevivable. The claim is advisory: revive itself does not require it,
	// because rejecting a late revive cannot un-spend anyone's item.
	ReviveClaimSource uuid.UUID
	ReviveClaimExpiry time.Time
}

// HasActiveReviveClaim reports whether an unexpired revive claim stands on this
// player at time t. An expired claim is simply not active — nothing needs to
// sweep it; the next claim attempt overwrites it.
func (p *LivePlayer) HasActiveReviveClaim(t time.Time) bool {
	return uuid.Equal(p.ReviveClaimSource, uuid.Nil) == false && t.Before(p.ReviveClaimExpiry)
}

var (
	ErrPlayerActionsLocked = errors.New("player actions locked")
	ErrPlayerActionsFull   = errors.New("player actions full")
	ErrPlayerActionsEmpty  = errors.New("player actions empty")
	ErrPlayerDead          = errors.New("player is dead")
)

func (p *LivePlayer) CanEnqueueAction() error {

	if p.Dead {
		return ErrPlayerDead
	}

	if p.ActionsLocked {
		return ErrPlayerActionsLocked
	}

	if len(p.Actions) >= p.MaxActionCount && p.MaxActionCount > 0 {
		return ErrPlayerActionsFull
	}

	return nil
}

func (p *LivePlayer) CanDequeueAction() error {

	if p.Dead {
		return ErrPlayerDead
	}

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
