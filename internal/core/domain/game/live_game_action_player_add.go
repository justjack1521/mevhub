package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
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

	game.SendChange(NewPlayerAddChange(a.UserID, a.PlayerID, a.PartySlot))

	return nil

}
