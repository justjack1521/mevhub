package memory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/dto"
	"mevhub/internal/domain/lobby"
	"strings"
)

const lobbySummaryKey = "lobby_instance"
const lobbyPartyKey = "lobby_party"

var (
	ErrLobbyInstanceNotFoundByKey = func(key string) error {
		return fmt.Errorf("lobby instance not found by key: %s", key)
	}
)

type LobbyInstanceRedisRepository struct {
	client *redis.Client
}

func NewLobbyInstanceRedisRepository(client *redis.Client) *LobbyInstanceRedisRepository {
	return &LobbyInstanceRedisRepository{client: client}
}

func (r *LobbyInstanceRedisRepository) QueryByID(ctx context.Context, id uuid.UUID) (*lobby.Instance, error) {

	var key = r.GenerateLobbyInstanceKey(id)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, lobby.ErrFailedQueryLobbyInstance(ErrLobbyInstanceNotFoundByKey(key))
	}

	cmd := r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	instance := &dto.LobbyInstanceRedis{}
	if err := cmd.Scan(instance); err != nil {
		return nil, err
	}

	return instance.ToEntity(), nil
}

func (r *LobbyInstanceRedisRepository) QueryByPartyID(ctx context.Context, party string) (*lobby.Instance, error) {
	result, err := r.client.Get(ctx, r.GenerateLobbyPartyKey(party)).Result()
	if err != nil {
		return nil, err
	}
	id, err := uuid.FromString(result)
	if err != nil {
		return nil, err
	}
	return r.QueryByID(ctx, id)
}

func (r *LobbyInstanceRedisRepository) Create(ctx context.Context, instance *lobby.Instance) error {
	var key = r.GenerateLobbyInstanceKey(instance.SysID)
	var result = &dto.LobbyInstanceRedis{
		LobbyID:            instance.SysID.String(),
		QuestID:            instance.QuestID.String(),
		HostID:             instance.HostID.String(),
		PartyID:            instance.PartyID,
		MinimumPlayerLevel: instance.MinimumPlayerLevel,
		RegisteredAt:       instance.RegisteredAt.Unix(),
	}
	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return lobby.ErrFailedCreateLobbyInstance(err)
	}
	r.client.Expire(ctx, key, lobby.KeepAliveTime)

	if err := r.client.Set(ctx, r.GenerateLobbyPartyKey(instance.PartyID), result.LobbyID, lobby.KeepAliveTime).Err(); err != nil {
		return lobby.ErrFailedCreateLobbyInstance(err)
	}

	return nil

}

func (r *LobbyInstanceRedisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var key = r.GenerateLobbyInstanceKey(id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return lobby.ErrFailedDeleteLobbyInstance(err)
	}
	return nil
}

func (r *LobbyInstanceRedisRepository) GenerateLobbyInstanceKey(id uuid.UUID) string {
	return strings.Join([]string{lobbySummaryKey, id.String()}, ":")
}

func (r *LobbyInstanceRedisRepository) GenerateLobbyPartyKey(party string) string {
	return strings.Join([]string{lobbyPartyKey, party}, ":")
}
