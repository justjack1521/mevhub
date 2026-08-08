package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

// evictExpiredDisconnectedPlayers removes any player whose reconnect grace
// period has lapsed, emitting the same PlayerRemoveChange an explicit remove
// action would. Every state calls this from Update so a turn-boundary stall on
// disconnected players is always bounded by the domain itself, even if the
// application layer's timeout eviction fails or is never delivered.
func evictExpiredDisconnectedPlayers(instance *game.LiveGameInstance, t time.Time) {
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			if !player.Disconnected || player.DisconnectTime.IsZero() {
				continue
			}
			if t.Sub(player.DisconnectTime) < game.DisconnectGracePeriod {
				continue
			}
			if err := party.RemovePlayer(player.PlayerID); err != nil {
				continue
			}
			instance.SendChange(NewPlayerRemoveChange(instance.InstanceID, player.UserID, player.PlayerID, party.PartyIndex, player.PartySlot))
		}
	}
}
