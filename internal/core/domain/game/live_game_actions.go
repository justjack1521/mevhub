package game

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"time"
)

var (
	ErrPartyIDNil          = errors.New("party id is nil")
	ErrPartyAlreadyInGame  = errors.New("party already added to game")
	ErrPlayerGameFull      = errors.New("live game is full")
	ErrPlayerAlreadyInGame = errors.New("player already added to game")
	ErrPlayerNotInGame     = errors.New("player not in game")
	ErrPlayerNotInParty    = errors.New("player not in party")
	ErrUserIDNil           = errors.New("user id is nil")
	ErrPlayerIDNil         = errors.New("player id is nil")
)

type Action interface {
	Perform(game *LiveGameInstance) error
}

type PartyAddAction struct {
	PartyID    uuid.UUID
	PartyIndex int
}

func (a *PartyAddAction) Perform(game *LiveGameInstance) error {

	if game.PartyExists(a.PartyID) {
		return ErrPartyAlreadyInGame
	}

	if uuid.Equal(a.PartyID, uuid.Nil) {
		return ErrPartyIDNil
	}

	game.Parties[a.PartyID] = &LiveParty{
		PartyID:        a.PartyID,
		PartyIndex:     a.PartyIndex,
		Players:        make(map[uuid.UUID]*LivePlayer),
		MaxPlayerCount: 100,
		LastAction:     time.Now().UTC(),
	}

	var change = PartyAddChange{
		PartyID:   a.PartyID,
		PartySlot: a.PartyIndex,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedAddPlayerToGame = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to add player %s: %w", player, err)
	}
)

type PlayerAddAction struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartyID   uuid.UUID
	PartySlot int
}

func (a *PlayerAddAction) validate(game *LiveGameInstance) error {

	party, err := game.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	if party.GetPlayerCount() == party.MaxPlayerCount {
		return ErrPlayerGameFull
	}

	if uuid.Equal(a.UserID, uuid.Nil) {
		return ErrUserIDNil
	}

	if uuid.Equal(a.PlayerID, uuid.Nil) {
		return ErrPlayerIDNil
	}

	if _, exists := party.Players[a.PlayerID]; exists {
		return ErrPlayerAlreadyInGame
	}
	return nil
}

func (a *PlayerAddAction) Perform(game *LiveGameInstance) error {

	if err := a.validate(game); err != nil {
		return ErrFailedAddPlayerToGame(a.PlayerID, err)
	}

	party, err := game.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	party.Players[a.PlayerID] = &LivePlayer{
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
	ErrFailedRemovePlayer = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to remove player %s: %w", id, err)
	}
)

type PlayerRemoveAction struct {
	InstanceID uuid.UUID
	UserID     uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerRemoveAction) Perform(game *LiveGameInstance) error {

	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedRemovePlayer(a.PlayerID, err)
	}

	if err := game.RemovePlayer(a.PlayerID); err != nil {
		return ErrFailedRemovePlayer(a.PlayerID, err)
	}

	var change = PlayerRemoveChange{
		UserID:    player.UserID,
		PlayerID:  player.PlayerID,
		PartySlot: player.PartySlot,
	}

	game.ChangeChannel <- change

	return nil

}

var (
	ErrFailedReadyPlayer = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to ready player %s: %w", player, err)
	}
)

type PlayerReadyAction struct {
	GameID   uuid.UUID
	PartyID  uuid.UUID
	PlayerID uuid.UUID
}

func (a *PlayerReadyAction) Perform(game *LiveGameInstance) error {

	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedReadyPlayer(a.PlayerID, err)
	}

	player.Ready = true
	var change = PlayerReadyChange{
		InstanceID: a.GameID,
		PartySlot:  player.PartySlot,
	}
	game.SendChange(change)
	return nil

}

var (
	ErrFailedDisconnectPlayer = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to disconnect player %s: %w", player, err)
	}
)

type PlayerDisconnectAction struct {
	InstanceID     uuid.UUID
	PlayerID       uuid.UUID
	DisconnectTime time.Time
}

func (a *PlayerDisconnectAction) Perform(game *LiveGameInstance) error {
	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDisconnectPlayer(a.PlayerID, err)
	}
	player.DisconnectTime = a.DisconnectTime
	return nil
}

type StateChangeAction struct {
	InstanceID uuid.UUID
	State      State
}

func (a *StateChangeAction) Perform(game *LiveGameInstance) error {

	fmt.Println("Change state to", reflect.TypeOf(a.State), " at ", time.Now().UTC().String())
	game.State = a.State

	var change = StateChange{
		InstanceID: a.InstanceID,
		State:      a.State,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedEnqueueAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to enqueue action for player %s: %w", player, err)
	}
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

	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedEnqueueAction(a.PlayerID, err)
	}

	var action = &PlayerAction{
		Target:     a.Target,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		ElementID:  a.ElementID,
	}

	if err := player.EnqueueAction(action); err != nil {
		return ErrFailedEnqueueAction(a.PlayerID, err)
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
	ErrFailedDequeueAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s: %w", player, err)
	}
)

type PlayerDequeueAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerDequeueAction) Perform(game *LiveGameInstance) error {

	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDequeueAction(a.PlayerID, err)
	}

	if err := player.DequeueAction(); err != nil {
		return ErrFailedDequeueAction(a.PlayerID, err)
	}

	var change = PlayerDequeueActionChange{
		InstanceID: a.InstanceID,
		PartySlot:  player.PartySlot,
	}

	game.SendChange(change)

	return nil

}

var (
	ErrFailedLockAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s: %w", player, err)
	}
	ErrPlayerUnableToLockAction = errors.New("player unable to dequeue action")
)

type PlayerLockAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerLockAction) Perform(game *LiveGameInstance) error {

	player, err := game.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedLockAction(a.PlayerID, err)
	}

	if player.ActionsLocked {
		return ErrFailedLockAction(a.PlayerID, ErrPlayerUnableToLockAction)
	}

	player.ActionLockIndex = game.GetActionLockedPlayerCount()
	player.ActionsLocked = true

	var change = PlayerLockActionChange{
		InstanceID:      a.InstanceID,
		PartySlot:       player.PartySlot,
		ActionLockIndex: player.ActionLockIndex,
	}

	game.SendChange(change)

	return nil

}
