package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

// evictExpiredDeadPlayers removes any player who has been dead longer than the
// game's DeadPlayerKickDuration without being revived, emitting the same
// PlayerRemoveChange an explicit remove action would. Every state calls this
// from Update, alongside the disconnect eviction it mirrors.
//
// A zero duration disables the kick entirely — the current configuration, since
// nothing populates the option yet — so dead players simply wait for a revive.
func evictExpiredDeadPlayers(instance *game.LiveGameInstance, t time.Time) {
	if instance.DeadPlayerKickDuration <= 0 {
		return
	}
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			if !player.Dead || player.DeathTime.IsZero() {
				continue
			}
			if t.Sub(player.DeathTime) < instance.DeadPlayerKickDuration {
				continue
			}
			if err := party.RemovePlayer(player.PlayerID); err != nil {
				continue
			}
			instance.SendChange(NewPlayerRemoveChange(instance.InstanceID, player.UserID, player.PlayerID, party.PartyIndex, player.PartySlot))
		}
	}
}
