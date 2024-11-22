package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrFailedRemovePlayer = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to remove player %s: %w", id, err)
	}
)

type PlayerRemoveAction struct {
	GameID   uuid.UUID
	PartyID  uuid.UUID
	UserID   uuid.UUID
	PlayerID uuid.UUID
}

func NewPlayerRemoveAction(instanceID uuid.UUID, partyID uuid.UUID, userID uuid.UUID, playerID uuid.UUID) *PlayerRemoveAction {
	return &PlayerRemoveAction{GameID: instanceID, PartyID: partyID, UserID: userID, PlayerID: playerID}
}

func (a *PlayerRemoveAction) Perform(game *LiveGameInstance) error {

	party, err := game.GetParty(a.PartyID)
	if err != nil {
		return ErrFailedRemovePlayer(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedRemovePlayer(a.PlayerID, err)
	}

	if err := game.RemovePlayer(a.PlayerID); err != nil {
		return ErrFailedRemovePlayer(a.PlayerID, err)
	}

	game.ChangeChannel <- NewPlayerRemoveChange(player.UserID, player.PlayerID, party.PartyIndex, player.PartySlot)

	return nil

}
