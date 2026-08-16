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

	// Idempotency: a duplicate ConnectedEvent for an already-connected player
	// must not re-broadcast a reconnect.
	if !player.Disconnected {
		return nil
	}

	player.Disconnected = false
	player.DisconnectTime = time.Time{}
	instance.SendChange(NewPlayerReconnectChange(instance.InstanceID, a.PlayerID, party.PartyIndex, player.PartySlot))
	// No state snapshot is pushed here any more. A reconnecting client re-syncs
	// itself by calling GameCatchUp with the last sequence it applied, which
	// replays what it missed through the handlers it already runs — rather than
	// through a reconciliation path exercised only on reconnect.
	// NewGameStateSyncChange and its marshaller are left intact and dormant.
	return nil
}
