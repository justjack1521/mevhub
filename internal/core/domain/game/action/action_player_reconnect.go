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
	PlayerID      uuid.UUID
	ReconnectTime time.Time
}

func NewPlayerReconnectAction(playerID uuid.UUID, reconnectTime time.Time) *PlayerReconnectAction {
	return &PlayerReconnectAction{PlayerID: playerID, ReconnectTime: reconnectTime}
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

	// Back mid enemy turn, they may still have the round to play out - and if the
	// queue confirm went out while they were away, they have not even started it.
	// The ready wait's timeout ran on while they were gone, so give them a fresh one.
	if state, ok := instance.State.(*EnemyTurnState); ok {
		state.restartReadyWait(a.ReconnectTime)
	}

	// Restate immediately rather than leaving the returning player to wait out
	// the heartbeat. Nothing changes state on a reconnect, so this is the only
	// prompt emission they would get — and because both turn boundaries stall
	// while anyone is disconnected, what they missed is bounded to this state,
	// all of which the restatement describes.
	restateGame(instance)

	return nil
}
