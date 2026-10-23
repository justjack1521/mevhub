package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

var (
	ErrFailedPlayerStamp = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to send stamp for player %s: %w", player, err)
	}
)

// PlayerStampAction relays a stamp to everyone in the game. Like chat, the party
// lookup only resolves the sender into the party/slot pair clients show it over,
// and dead or disconnected players may still send one. The one rule of its own
// is the cooldown: a stamp inside StampCooldown of the player's last is dropped
// without a reply - a refused tap is not a failure worth reporting.
type PlayerStampAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	StampID    uuid.UUID
	StampTime  time.Time
}

func NewPlayerStampAction(instanceID uuid.UUID, playerID uuid.UUID, stampID uuid.UUID, stampTime time.Time) *PlayerStampAction {
	return &PlayerStampAction{InstanceID: instanceID, PlayerID: playerID, StampID: stampID, StampTime: stampTime}
}

func (a *PlayerStampAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerStamp(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerStamp(a.PlayerID, err)
	}

	if player.CanStamp(a.StampTime) == false {
		return nil
	}

	player.LastStampTime = a.StampTime

	instance.SendChange(NewPlayerStampChange(instance.InstanceID, party.PartyIndex, player.PartySlot, a.StampID))
	return nil

}
