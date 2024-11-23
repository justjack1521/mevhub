package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

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

func NewPlayerReadyAction(gameID uuid.UUID, partyID uuid.UUID, playerID uuid.UUID) *PlayerReadyAction {
	return &PlayerReadyAction{GameID: gameID, PartyID: partyID, PlayerID: playerID}
}

func (a *PlayerReadyAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetParty(a.PartyID)
	if err != nil {
		return ErrFailedReadyPlayer(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedReadyPlayer(a.PlayerID, err)
	}

	player.Ready = true

	instance.SendChange(NewPlayerReadyChange(a.GameID, party.PartyIndex, player.PartySlot))
	return nil

}
