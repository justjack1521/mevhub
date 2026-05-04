package memory

import (
	"context"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/core/domain/game"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

const (
	gamePlayerLoadoutKey = "game_player_loadout"
	gamePlayerLoadoutTTL = time.Hour * 3
)

type GamePlayerLoadoutRepository struct {
	client     *redis.Client
	serialiser serial.GamePlayerLoadoutSerialiser
}

func NewGamePlayerLoadoutRepository(client *redis.Client, serialiser serial.GamePlayerLoadoutSerialiser) *GamePlayerLoadoutRepository {
	return &GamePlayerLoadoutRepository{client: client, serialiser: serialiser}
}

func (r *GamePlayerLoadoutRepository) Query(ctx context.Context, player uuid.UUID, index int) (game.PlayerLoadout, error) {
	value, err := r.client.Get(ctx, r.Key(player, index)).Result()
	if err != nil {
		return game.PlayerLoadout{}, err
	}
	return r.serialiser.Unmarshall([]byte(value))
}

func (r *GamePlayerLoadoutRepository) Create(ctx context.Context, player uuid.UUID, index int, loadout game.PlayerLoadout) error {
	result, err := r.serialiser.Marshall(loadout)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.Key(player, index), result, gamePlayerLoadoutTTL).Err()
}

func (r *GamePlayerLoadoutRepository) Delete(ctx context.Context, player uuid.UUID, index int) error {
	return r.client.Del(ctx, r.Key(player, index)).Err()
}

func (r *GamePlayerLoadoutRepository) Key(player uuid.UUID, index int) string {
	return strings.Join([]string{serviceKey, gamePlayerLoadoutKey, player.String(), strconv.Itoa(index)}, ":")
}
