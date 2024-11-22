package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

var (
	ErrFailedDequeueAction = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to dequeue action for player %s: %w", player, err)
	}
)

type PlayerDequeueAction struct {
	InstanceID uuid.UUID
	PartyID    uuid.UUID
	PlayerID   uuid.UUID
}

func NewPlayerDequeueAction(instanceID uuid.UUID, partyID uuid.UUID, playerID uuid.UUID) *PlayerDequeueAction {
	return &PlayerDequeueAction{InstanceID: instanceID, PartyID: partyID, PlayerID: playerID}
}

func (a *PlayerDequeueAction) Perform(game *game.LiveGameInstance) error {

	party, err := game.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDequeueAction(a.PlayerID, err)
	}

	if err := player.DequeueAction(); err != nil {
		return ErrFailedDequeueAction(a.PlayerID, err)
	}

	game.SendChange(NewPlayerDequeueActionChange(a.InstanceID, party.PartyIndex, player.PartySlot))

	return nil

}
