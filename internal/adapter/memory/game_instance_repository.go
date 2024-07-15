package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/domain/game"
	"strings"
	"time"
)

const (
	gameInstanceKey = "game_instance"
	gameInstanceTTL = time.Hour * 3
)

type GameInstanceRepository struct {
	client     *redis.Client
	serialiser serial.GameInstanceSerialiser
}

func NewGameInstanceRepository(client *redis.Client, serialiser serial.GameInstanceSerialiser) *GameInstanceRepository {
	return &GameInstanceRepository{client: client, serialiser: serialiser}
}

func (r *GameInstanceRepository) Get(ctx context.Context, id uuid.UUID) (*game.Instance, error) {
	value, err := r.client.Get(ctx, r.Key(id)).Result()
	if err != nil {
		return nil, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GameInstanceRepository) Create(ctx context.Context, instance *game.Instance) error {
	result, err := r.serialiser.Marshall(instance)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(instance.SysID), result, gameInstanceTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GameInstanceRepository) Key(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, gameInstanceKey, id.String()}, ":")
}
