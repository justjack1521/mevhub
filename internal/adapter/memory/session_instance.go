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

	exists, err := r.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	var key = r.GenerateSessionKey(id)

	if exists == false {
		return nil, ErrSessionInstanceNotFoundByKey(key)
	}

	var cmd = r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var result = &dto.SessionInstanceRedis{}
	if err := cmd.Scan(result); err != nil {
		return nil, err
	}
	r.client.Expire(ctx, key, time.Minute*30)
	return result.ToEntity(), nil
}

func (r *SessionInstanceRedisRepository) Create(ctx context.Context, instance *session.Instance) error {

	result, err := r.InstanceToTransfer(instance)
	if err != nil {
		return err
	}

	var key = r.GenerateSessionKey(instance.UserID)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return err
	}
	r.client.Expire(ctx, key, time.Minute*120)
	return nil
}

func (r *SessionInstanceRedisRepository) Update(ctx context.Context, instance *session.Instance) error {
	result, err := r.InstanceToTransfer(instance)
	if err != nil {
		return err
	}
	var key = r.GenerateSessionKey(instance.UserID)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return err
	}
	r.client.Expire(ctx, key, time.Minute*30)
	return nil
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
		UserID:    instance.UserID.String(),
		PlayerID:  instance.PlayerID.String(),
		DeckIndex: instance.DeckIndex,
		LobbyID:   instance.LobbyID.String(),
		PartySlot: instance.PartySlot,
	}, nil
}

func (r *SessionInstanceRedisRepository) GenerateSessionKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, sessionKey, id.String()}, sessionKeySeparator)
}
