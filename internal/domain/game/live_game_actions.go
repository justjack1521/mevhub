package game

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrPlayerGameFull      = errors.New("live game is full")
	ErrPlayerAlreadyInGame = errors.New("player already added to game")
	ErrPlayerNotInGame     = errors.New("player not in game")
	ErrUserIDNil           = errors.New("user id is nil")
	ErrPlayerIDNil         = errors.New("player id is nil")
)

type Action interface {
	Perform(game *LiveGameInstance) error
}

var (
	ErrFailedAddPlayerToGame = func(player uuid.UUID, instance uuid.UUID, err error) error {
		return fmt.Errorf("failed to add player %s to live game %s: %w", player, instance, err)
	}
)

type PlayerAddAction struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

func (a *PlayerAddAction) validate(game *LiveGameInstance) error {
	if len(game.Players) == game.MaxPlayerCount {
		return ErrPlayerGameFull
	}

	if uuid.Equal(a.UserID, uuid.Nil) {
		return ErrUserIDNil
	}

	if uuid.Equal(a.PlayerID, uuid.Nil) {
		return ErrPlayerIDNil
	}

	if _, exists := game.Players[a.PlayerID]; exists {
		return ErrPlayerAlreadyInGame
	}
	return nil
}

func (a *PlayerAddAction) Perform(game *LiveGameInstance) error {

	if err := a.validate(game); err != nil {
		return ErrFailedAddPlayerToGame(a.PlayerID, game.InstanceID, err)
	}

	game.Players[a.PlayerID] = &LivePlayer{
		UserID:    a.UserID,
		PlayerID:  a.PlayerID,
		PartySlot: a.PartySlot,
	}

	var change = PlayerAddChange{
		UserID:    a.UserID,
		PlayerID:  a.PlayerID,
		PartySlot: a.PartySlot,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedReadyPlayer = func(player uuid.UUID, instance uuid.UUID, err error) error {
		return fmt.Errorf("failed to ready player %s in live game %s: %w", player, instance, err)
	}
)

type PlayerReadyAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerReadyAction) Perform(game *LiveGameInstance) error {

	player, valid := game.Players[a.PlayerID]
	if valid == false {
		return ErrFailedReadyPlayer(a.PlayerID, game.InstanceID, ErrPlayerNotInGame)
	}

	player.Ready = true
	var change = PlayerReadyChange{
		InstanceID: a.InstanceID,
		PartySlot:  player.PartySlot,
	}
	game.SendChange(change)
	return nil

}

type StateChangeAction struct {
	InstanceID uuid.UUID
	State      State
}

func (a *StateChangeAction) Perform(game *LiveGameInstance) error {

	game.State = a.State

	var change = StateChange{
		InstanceID: a.InstanceID,
		State:      game.State,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedEnqueueAction = func(player uuid.UUID, instance uuid.UUID, err error) error {
		return fmt.Errorf("failed to enqueue action for player %s in live game %s: %w", player, instance, err)
	}
	ErrPlayerUnableToEnqueueAction = errors.New("player unable to enqueue action")
)

type PlayerEnqueueAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	Target     int
	ActionType PlayerActionType
	SlotIndex  int
	ElementID  uuid.UUID
}

func (a *PlayerEnqueueAction) Perform(game *LiveGameInstance) error {

	player, valid := game.Players[a.PlayerID]
	if valid == false {
		return ErrFailedEnqueueAction(a.PlayerID, game.InstanceID, ErrPlayerNotInGame)
	}

	var action = &PlayerAction{
		Target:     a.Target,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		ElementID:  a.ElementID,
	}

	if player.EnqueueAction(action) == false {
		return ErrFailedEnqueueAction(a.PlayerID, game.InstanceID, ErrPlayerUnableToEnqueueAction)
	}

	var change = PlayerEnqueueActionChange{
		InstanceID: game.InstanceID,
		PartySlot:  player.PartySlot,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		Target:     a.Target,
		ElementID:  a.ElementID,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedDequeueAction = func(player uuid.UUID, instance uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s in live game %s: %w", player, instance, err)
	}
	ErrPlayerUnableToDequeueAction = errors.New("player unable to dequeue action")
)

type PlayerDequeueAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerDequeueAction) Perform(game *LiveGameInstance) error {

	player, exists := game.Players[a.PlayerID]
	if exists == false {
		return ErrFailedDequeueAction(a.PlayerID, game.InstanceID, ErrPlayerNotInGame)
	}

	if player.DequeueAction() == false {
		return ErrFailedDequeueAction(a.PlayerID, game.InstanceID, ErrPlayerUnableToDequeueAction)
	}

	var change = PlayerDequeueActionChange{
		InstanceID: a.InstanceID,
		PartySlot:  player.PartySlot,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedLockAction = func(player uuid.UUID, instance uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s in live game %s: %w", player, instance, err)
	}
	ErrPlayerUnableToLockAction = errors.New("player unable to dequeue action")
)

type PlayerLockAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerLockAction) Perform(game *LiveGameInstance) error {

	player, exists := game.Players[a.PlayerID]
	if exists == false {
		return ErrFailedLockAction(a.PlayerID, game.InstanceID, ErrPlayerNotInGame)
	}

	if player.ActionsLocked {
		return ErrFailedLockAction(a.PlayerID, game.InstanceID, ErrPlayerUnableToLockAction)
	}

	player.ActionsLocked = true
	player.ActionLockIndex = game.GetActionLockedPlayerCount()

	var change = PlayerLockActionChange{
		InstanceID:      a.InstanceID,
		PartySlot:       player.PartySlot,
		ActionLockIndex: player.ActionLockIndex,
	}

	game.SendChange(change)

	return nil

}
