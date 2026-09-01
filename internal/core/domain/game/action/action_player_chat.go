package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

var (
	ErrFailedPlayerChat = func(player uuid.UUID, err error) error {
		return fmt.Errorf("failed to send chat for player %s: %w", player, err)
	}
)

// PlayerChatAction relays a chat message to everyone in the game. It mutates
// no game state — the party lookup exists only to resolve the sender into the
// party/slot pair clients use to attribute the message. Dead and disconnected
// players may still chat; only a player no longer in the instance is refused.
type PlayerChatAction struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
	Message    string
}

func NewPlayerChatAction(instanceID uuid.UUID, playerID uuid.UUID, message string) *PlayerChatAction {
	return &PlayerChatAction{InstanceID: instanceID, PlayerID: playerID, Message: message}
}

func (a *PlayerChatAction) Perform(instance *game.LiveGameInstance) error {

	party, err := instance.GetPartyForPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerChat(a.PlayerID, err)
	}

	player, err := party.GetPlayer(a.PlayerID)
	if err != nil {
		return ErrFailedPlayerChat(a.PlayerID, err)
	}

	instance.SendChange(NewPlayerChatChange(instance.InstanceID, party.PartyIndex, player.PartySlot, a.Message))
	return nil

}
