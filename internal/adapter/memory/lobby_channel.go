package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"strings"
)

// Listeners are stored as a hash of user ID -> player ID. Both are needed: the
// client notification is routed on the player ID header and a nil one is
// discarded downstream. The key name differs from the previous set-based one so
// live set keys are not hit with WRONGTYPE while they age out.
const lobbyChannelKey = "lobby_notification_listener"

type LobbyChannelRepository struct {
	client *redis.Client
}

func NewLobbyChannelRepository(client *redis.Client) *LobbyChannelRepository {
	return &LobbyChannelRepository{client: client}
}

func (r *LobbyChannelRepository) QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]lobby.NotificationListener, error) {
	var key = r.GenerateKeyForLobby(id)
	results, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var members = make([]lobby.NotificationListener, 0, len(results))
	for user, player := range results {
		u, err := uuid.FromString(user)
		if err != nil {
			continue
		}
		p, err := uuid.FromString(player)
		if err != nil {
			continue
		}
		members = append(members, lobby.NotificationListener{LobbyID: id, UserID: u, PlayerID: p})
	}
	return members, nil
}

func (r *LobbyChannelRepository) CreateListener(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID) error {
	var key = r.GenerateKeyForLobby(id)
	if err := r.client.HSet(ctx, key, user.String(), player.String()).Err(); err != nil {
		return err
	}
	r.client.Expire(ctx, key, lobby.KeepAliveTime)
	return nil
}

func (r *LobbyChannelRepository) DeleteListener(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	var key = r.GenerateKeyForLobby(id)
	if err := r.client.HDel(ctx, key, user.String()).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LobbyChannelRepository) DeleteAll(ctx context.Context, id uuid.UUID) error {
	var key = r.GenerateKeyForLobby(id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LobbyChannelRepository) GenerateKeyForLobby(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, lobbyChannelKey, id.String()}, ":")
}
