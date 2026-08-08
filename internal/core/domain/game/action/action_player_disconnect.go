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
	PlayerID       uuid.UUID
	DisconnectTime time.Time
}

func NewPlayerDisconnectAction(instanceID uuid.UUID, playerID uuid.UUID, disconnectTime time.Time) *PlayerDisconnectAction {
	return &PlayerDisconnectAction{InstanceID: instanceID, PlayerID: playerID, DisconnectTime: disconnectTime}
}

func (a *PlayerDisconnectAction) Perform(instance *game.LiveGameInstance) error {

	// The session no longer carries a party reference once the game starts,
	// so resolve the party from the player, as reconnect/remove already do.
	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDisconnectPlayer(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedDisconnectPlayer(a.PlayerID, err)
	}

	player.Disconnected = true
	player.DisconnectTime = a.DisconnectTime
	instance.SendChange(NewPlayerDisconnectChange(instance.InstanceID, a.PlayerID, party.PartyIndex, player.PartySlot))
	return nil
}
