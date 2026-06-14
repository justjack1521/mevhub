package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

var (
	ErrFailedReconnectPlayer = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to reconnect player %s: %w", player, err)
	}
)

type PlayerReconnectAction struct {
	PlayerID uuid.UUID
}

func NewPlayerReconnectAction(playerID uuid.UUID) *PlayerReconnectAction {
	return &PlayerReconnectAction{PlayerID: playerID}
}

func (a *PlayerReconnectAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedReconnectPlayer(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedReconnectPlayer(a.PlayerID, err)
	}

	player.Disconnected = false
	player.DisconnectTime = time.Time{}
	instance.ChangeChannel <- NewPlayerReconnectChange(instance.InstanceID, a.PlayerID, party.PartyIndex, player.PartySlot)
	return nil
}
