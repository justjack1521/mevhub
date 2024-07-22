package game

import uuid "github.com/satori/go.uuid"

type Action interface {
	Perform(game *LiveGameInstance)
}

type PlayerAddAction struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

func (a *PlayerAddAction) validate(game *LiveGameInstance) bool {
	if len(game.Players) == game.MaxPlayerCount {
		return false
	}

	if uuid.Equal(a.UserID, uuid.Nil) {
		return false
	}

	if uuid.Equal(a.PlayerID, uuid.Nil) {
		return false
	}

	if _, exists := game.Players[a.PlayerID]; exists {
		return false
	}
	return true
}

func (a *PlayerAddAction) Perform(game *LiveGameInstance) {

	if a.validate(game) == false {
		return
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

}

type PlayerReadyAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerReadyAction) Perform(game *LiveGameInstance) {
	if player, valid := game.Players[a.PlayerID]; valid {
		player.Ready = true
		var change = PlayerReadyChange{
			InstanceID: a.InstanceID,
			PartySlot:  player.PartySlot,
		}
		game.SendChange(change)
	}
}

type StateChangeAction struct {
	InstanceID uuid.UUID
	State      State
}

func (a *StateChangeAction) Perform(game *LiveGameInstance) {

	game.State = a.State

	var change = GameStateChange{
		InstanceID: a.InstanceID,
	}

	game.SendChange(change)

}

type PlayerEnqueueAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	Target     int
	ActionType PlayerActionType
	SlotIndex  int
	ElementID  uuid.UUID
}

func (a *PlayerEnqueueAction) Perform(game *LiveGameInstance) {

	player, exists := game.Players[a.PlayerID]
	if exists == false {
		return
	}

	var action = &PlayerAction{
		Target:     a.Target,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		ElementID:  a.ElementID,
	}

	if player.EnqueueAction(action) == false {
		return
	}

	var change = PlayerEnqueueActionChange{
		Change:     nil,
		InstanceID: game.InstanceID,
		PartySlot:  player.PartySlot,
		ActionType: a.ActionType,
		SlotIndex:  a.SlotIndex,
		Target:     a.Target,
		ElementID:  a.ElementID,
	}

	game.SendChange(change)

}

type PlayerDequeueAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerDequeueAction) Perform(game *LiveGameInstance) {

	player, exists := game.Players[a.PlayerID]
	if exists == false {
		return
	}

	if player.DequeueAction() == false {
		return
	}

	var change = PlayerDequeueActionChange{
		InstanceID: a.InstanceID,
		PartySlot:  player.PartySlot,
	}

	game.SendChange(change)

}

type PlayerLockAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}

func (a *PlayerLockAction) Perform(game *LiveGameInstance) {
	player, exists := game.Players[a.PlayerID]
	if exists == false {
		return
	}

	if player.ActionsLocked {
		return
	}

	player.ActionsLocked = true

	var change = PlayerLockActionChange{
		InstanceID: a.InstanceID,
		PartySlot:  player.PartySlot,
	}

	game.SendChange(change)

}
