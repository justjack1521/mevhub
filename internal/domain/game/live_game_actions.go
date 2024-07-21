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

func (a *PlayerAddAction) Perform(game *LiveGameInstance) {

	game.Mu.Lock()
	defer game.Mu.Unlock()

	game.Players[a.PlayerID] = &LivePlayer{
		UserID:    a.UserID,
		PlayerID:  a.PlayerID,
		PartySlot: a.PartySlot,
	}

	var change = PlayerAddChange{
		Change:    nil,
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
	game.Mu.Lock()
	defer game.Mu.Unlock()

	if player, valid := game.Players[a.PlayerID]; valid {
		player.Ready = true
		var change = PlayerReadyChange{
			InstanceID: a.InstanceID,
			PartySlot:  player.PartySlot,
		}
		game.SendChange(change)
	}

}
