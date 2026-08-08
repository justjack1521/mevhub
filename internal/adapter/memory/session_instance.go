package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/memory/dto"
	"mevhub/internal/core/domain/session"
	"strings"
	"time"
)

var (
	ErrSessionInstanceNotFoundByKey = func(key string) error {
		return fmt.Errorf("session instance not found by key: %s", key)
	}
)

const sessionKey = "lobby_session"
const sessionKeySeparator = ":"

// sessionTTL is the single session lifetime. It is refreshed on every write;
// reads deliberately do not touch it (a pure read must never shorten — or
// extend — a key's remaining life).
const sessionTTL = time.Minute * 120

type SessionInstanceRedisRepository struct {
	client *redis.Client
}

func NewLobbySessionRedisRepository(client *redis.Client) *SessionInstanceRedisRepository {
	return &SessionInstanceRedisRepository{client: client}
}

func (r *SessionInstanceRedisRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var key = r.GenerateSessionKey(id)
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *SessionInstanceRedisRepository) QueryByID(ctx context.Context, id uuid.UUID) (*session.Instance, error) {

	var key = r.GenerateSessionKey(id)

	// HGETALL on a missing key returns an empty map, not an error; a separate
	// EXISTS check races key expiry and previously let a zero-value session
	// escape as valid.
	var cmd = r.client.HGetAll(ctx, key)
	fields, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, ErrSessionInstanceNotFoundByKey(key)
	}
	var result = &dto.SessionInstanceRedis{}
	if err := cmd.Scan(result); err != nil {
		return nil, err
	}
	return result.ToEntity(), nil
}

func (r *SessionInstanceRedisRepository) Create(ctx context.Context, instance *session.Instance) error {
	result, err := r.InstanceToTransfer(instance)
	if err != nil {
		return err
	}
	return r.writeFields(ctx, instance.UserID, result.ToMapStringInterface())
}

func (r *SessionInstanceRedisRepository) Update(ctx context.Context, instance *session.Instance) error {
	result, err := r.InstanceToTransfer(instance)
	if err != nil {
		return err
	}
	return r.writeFields(ctx, instance.UserID, result.ToMapStringInterface())
}

func (r *SessionInstanceRedisRepository) UpdateConnectionState(ctx context.Context, instance *session.Instance) error {
	result, err := r.InstanceToTransfer(instance)
	if err != nil {
		return err
	}
	return r.writeFields(ctx, instance.UserID, result.ConnectionStateMap())
}

// writeFields applies an HSET and TTL refresh atomically. A bare HSET followed
// by a best-effort EXPIRE could leave a fresh key with no TTL at all — a
// permanently wedged session if the second round trip was lost.
func (r *SessionInstanceRedisRepository) writeFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	var key = r.GenerateSessionKey(id)
	pipe := r.client.TxPipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, sessionTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *SessionInstanceRedisRepository) Delete(ctx context.Context, instance *session.Instance) error {
	if err := r.client.Del(ctx, r.GenerateSessionKey(instance.UserID)).Err(); err != nil {
		return err
	}
	return nil
}

func (r *SessionInstanceRedisRepository) InstanceToTransfer(instance *session.Instance) (dto.SessionInstanceRedis, error) {
	if instance == nil {
		return dto.SessionInstanceRedis{}, errors.New("session instance is nil")
	}
	return dto.SessionInstanceRedis{
		UserID:              instance.UserID.String(),
		PlayerID:            instance.PlayerID.String(),
		DeckIndex:           instance.DeckIndex,
		LobbyID:             instance.LobbyID.String(),
		GameID:              instance.GameID.String(),
		PartySlot:           instance.PartySlot,
		DisconnectSessionID: instance.DisconnectSessionID.String(),
		LastConnEventAt:     instance.LastConnEventAt,
		DisconnectedAt:      instance.DisconnectedAt,
		CurrentSessionID:    instance.CurrentSessionID.String(),
	}, nil
}

func (r *SessionInstanceRedisRepository) GenerateSessionKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, sessionKey, id.String()}, sessionKeySeparator)
}
