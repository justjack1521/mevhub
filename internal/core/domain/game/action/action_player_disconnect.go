package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

var (
	ErrFailedDisconnectPlayer = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to disconnect player %s: %w", player, err)
	}
)

type PlayerDisconnectAction struct {
	InstanceID     uuid.UUID
	PartyID        uuid.UUID
	PlayerID       uuid.UUID
	DisconnectTime time.Time
}

func NewPlayerDisconnectAction(instanceID uuid.UUID, partyID uuid.UUID, playerID uuid.UUID, disconnectTime time.Time) *PlayerDisconnectAction {
	return &PlayerDisconnectAction{InstanceID: instanceID, PartyID: partyID, PlayerID: playerID, DisconnectTime: disconnectTime}
}

func (a *PlayerDisconnectAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetParty(a.PartyID)
	if err != nil {
		return err
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDisconnectPlayer(a.PlayerID, err)
	}

	player.DisconnectTime = a.DisconnectTime
	return nil
}
