package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"strings"
)

const lobbyChannelKey = "lobby_channel"

type LobbyChannelRepository struct {
	client *redis.Client
}

func NewLobbyChannelRepository(client *redis.Client) *LobbyChannelRepository {
	return &LobbyChannelRepository{client: client}
}

func (r *LobbyChannelRepository) QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]lobby.NotificationListener, error) {
	var key = r.GenerateKeyForLobby(id)
	results, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var members = make([]lobby.NotificationListener, 0)
	for _, value := range results {
		member, err := uuid.FromString(value)
		if err != nil {
			continue
		}
		members = append(members, lobby.NotificationListener{UserID: member})
	}
	return members, nil
}

func (r *LobbyChannelRepository) Create(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	var key = r.GenerateKeyForLobby(id)
	if err := r.client.SAdd(ctx, key, user.String()).Err(); err != nil {
		return err
	}
	r.client.Expire(ctx, key, lobby.KeepAliveTime)
	return nil
}

func (r *LobbyChannelRepository) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	var key = r.GenerateKeyForLobby(id)
	if err := r.client.SRem(ctx, key, user.String()).Err(); err != nil {
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
