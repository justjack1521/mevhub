package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/memory/dto"
	"mevhub/internal/core/domain/game"
	"strconv"
	"strings"
	"time"
)

const (
	gamePartyKey          = "game_party"
	gamePartyKeySeparator = ":"
	gamePartyTTL          = time.Hour * 3
)

type GamePartyRepository struct {
	client *redis.Client
}

func NewGamePartyRepository(client *redis.Client) *GamePartyRepository {
	return &GamePartyRepository{client: client}
}

func (r *GamePartyRepository) Create(ctx context.Context, id uuid.UUID, party *game.Party) error {

	var result = &dto.GamePartyRedis{
		SysID:     party.SysID.String(),
		PartyID:   party.PartyID,
		Index:     party.Index,
		PartyName: party.PartyName,
	}

	var key = r.Key(id, party.Index)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return err
	}

	r.client.Expire(ctx, key, gamePartyTTL)

	return nil
}

func (r *GamePartyRepository) Delete(ctx context.Context, id uuid.UUID, slot int) error {
	var key = r.Key(id, slot)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GamePartyRepository) DeleteAll(ctx context.Context, id uuid.UUID) error {
	keys, err := r.client.Keys(ctx, r.GameKey(id)).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return err
	}

	return nil
}

func (r *GamePartyRepository) Get(ctx context.Context, id uuid.UUID) (*game.Party, error) {
	//TODO implement me
	panic("implement me")
}

func (r *GamePartyRepository) query(ctx context.Context, key string) (*game.Party, error) {
	var cmd = r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var result = &dto.GamePartyRedis{}
	if err := cmd.Scan(result); err != nil {
		return nil, err
	}
	return result.ToEntity(), nil
}

func (r *GamePartyRepository) Query(ctx context.Context, id uuid.UUID, slot int) (*game.Party, error) {
	var key = r.Key(id, slot)
	return r.query(ctx, key)
}

func (r *GamePartyRepository) QueryAll(ctx context.Context, id uuid.UUID) ([]*game.Party, error) {
	keys, err := r.client.Keys(ctx, r.GameKey(id)).Result()
	if err != nil {
		return nil, err
	}
	var participants = make([]*game.Party, len(keys))
	for index, key := range keys {
		participant, err := r.query(ctx, key)
		if err != nil {
			return nil, err
		}
		participants[index] = participant
	}
	return participants, nil
}

func (r *GamePartyRepository) GameKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, gamePartyKey, id.String(), "*"}, gamePartyKeySeparator)
}

func (r *GamePartyRepository) Key(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, gamePartyKey, id.String(), strconv.Itoa(slot)}, gamePartyKeySeparator)
}
