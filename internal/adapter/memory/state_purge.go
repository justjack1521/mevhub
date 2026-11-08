package memory

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// purgeScanBatch is the SCAN hint and the size of each UNLINK. Small enough
// that no single round trip holds Redis up, large enough that a few thousand
// stale keys clear in a handful of calls.
const purgeScanBatch = 500

// PurgeServiceState deletes every key this service owns, so a restart begins
// from nothing. Run once at startup, before anything is served.
//
// Nothing under the service prefix outlives the process that wrote it: live
// games exist only in that process's memory, so after a restart every session,
// lobby, search entry, queue and game record left in Redis points at something
// that no longer exists. Left in place, fresh clients inherit it - SessionCreate
// carries a stale session's game forward, and search lists lobbies with no live
// channel behind them. The loadout and player summary caches go too; they are
// refetched on demand.
//
// Scoped to the service prefix, so a Redis shared with other services keeps
// their keys. UNLINK frees memory off the main thread.
func PurgeServiceState(ctx context.Context, client *redis.Client) (int, error) {
	var cursor uint64
	var purged = 0
	for {
		keys, next, err := client.Scan(ctx, cursor, serviceKey+":*", purgeScanBatch).Result()
		if err != nil {
			return purged, err
		}
		if len(keys) > 0 {
			if err := client.Unlink(ctx, keys...).Err(); err != nil {
				return purged, err
			}
			purged += len(keys)
		}
		cursor = next
		if cursor == 0 {
			return purged, nil
		}
	}
}
