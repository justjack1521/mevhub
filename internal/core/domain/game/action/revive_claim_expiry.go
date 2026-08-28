package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

// expireLapsedReviveClaims clears any revive claim whose window has passed
// without the revive landing, broadcasting the expiry so clients can drop
// their "being revived" indicator. Every state calls this from Update,
// alongside the death kick sweep it mirrors.
//
// This is the only place a lapsed claim is observed: claims that end in a
// revive are cleared (and signalled) by the revive itself, and a kicked
// player's claim vanishes with the player, covered by PlayerRemoveChange.
// Clearing the fields is also what makes the broadcast fire exactly once.
func expireLapsedReviveClaims(instance *game.LiveGameInstance, t time.Time) {
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			if uuid.Equal(player.ReviveClaimSource, uuid.Nil) || t.Before(player.ReviveClaimExpiry) {
				continue
			}
			var source = player.ReviveClaimSource
			player.ReviveClaimSource = uuid.Nil
			player.ReviveClaimExpiry = time.Time{}
			instance.SendChange(NewPlayerReviveClaimExpireChange(instance.InstanceID, player.PlayerID, source, party.PartyIndex, player.PartySlot))
		}
	}
}
