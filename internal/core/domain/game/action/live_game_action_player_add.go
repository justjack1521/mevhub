package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

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

func NewPlayerAddAction(userID uuid.UUID, playerID uuid.UUID, partyID uuid.UUID, partySlot int) *PlayerAddAction {
	return &PlayerAddAction{UserID: userID, PlayerID: playerID, PartyID: partyID, PartySlot: partySlot}
}

func (a *PlayerAddAction) validate(instance *game.LiveGameInstance) error {

	party, err := instance.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	if party.GetPlayerCount() == party.MaxPlayerCount {
		return game.ErrPlayerGameFull
	}

	if uuid.Equal(a.UserID, uuid.Nil) {
		return game.ErrUserIDNil
	}

	if uuid.Equal(a.PlayerID, uuid.Nil) {
		return game.ErrPlayerIDNil
	}

	if _, exists := party.Players[a.PlayerID]; exists {
		return game.ErrPlayerAlreadyInGame
	}
	return nil
}

func (a *PlayerAddAction) Perform(instance *game.LiveGameInstance) error {

	if err := a.validate(instance); err != nil {
		return ErrFailedAddPlayerToGame(a.PlayerID, err)
	}

	party, err := instance.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	party.Players[a.PlayerID] = &game.LivePlayer{
		UserID:    a.UserID,
		PlayerID:  a.PlayerID,
		PartySlot: a.PartySlot,
	}

	instance.SendChange(NewPlayerAddChange(a.UserID, a.PlayerID, a.PartySlot))

	return nil

}
