package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/core/domain/game"
	"strconv"
	"strings"
	"time"
)

const (
	gamePlayerKey = "game_player"
	gamePlayerTTL = time.Hour * 3
)

type GamePlayerRepository struct {
	client     *redis.Client
	serialiser serial.GamePlayerSerialiser
}

func NewGamePlayerRepository(client *redis.Client, serialiser serial.GamePlayerSerialiser) *GamePlayerRepository {
	return &GamePlayerRepository{client: client, serialiser: serialiser}
}

func (r *GamePlayerRepository) Query(ctx context.Context, id uuid.UUID, slot int) (game.Player, error) {
	return r.query(ctx, r.Key(id, slot))
}

func (r *GamePlayerRepository) QueryAll(ctx context.Context, id uuid.UUID) ([]game.Player, error) {
	keys, err := r.client.Keys(ctx, r.GameKey(id)).Result()
	if err != nil {
		return nil, err
	}
	var participants = make([]game.Player, len(keys))
	for index, key := range keys {
		participant, err := r.query(ctx, key)
		if err != nil {
			return nil, err
		}
		participants[index] = participant
	}
	return participants, nil
}

func (r *GamePlayerRepository) Create(ctx context.Context, id uuid.UUID, slot int, participant game.Player) error {
	result, err := r.serialiser.Marshall(participant)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(id, slot), result, gamePlayerTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GamePlayerRepository) query(ctx context.Context, key string) (game.Player, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return game.Player{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return game.Player{}, err
	}
	return result, nil
}

func (r *GamePlayerRepository) GameKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, gamePlayerKey, id.String(), "*"}, ":")
}

func (r *GamePlayerRepository) Key(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, gamePlayerKey, id.String(), strconv.Itoa(slot)}, ":")
}
