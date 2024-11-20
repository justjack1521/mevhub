package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/core/domain/lobby"
	"strings"
	"time"
)

const (
	lobbySummaryKey = "lobby_summary"
	lobbySummaryTTL = time.Minute * 180
)

type LobbySummaryRepository struct {
	client     *redis.Client
	serialiser serial.LobbySummarySerialiser
}

func NewLobbySummaryRepository(client *redis.Client, serialiser serial.LobbySummarySerialiser) *LobbySummaryRepository {
	return &LobbySummaryRepository{client: client, serialiser: serialiser}
}

func (r *LobbySummaryRepository) Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error) {
	value, err := r.client.Get(ctx, r.Key(id)).Result()
	if err != nil {
		return lobby.Summary{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return lobby.Summary{}, err
	}
	return result, nil
}

func (r *LobbySummaryRepository) Create(ctx context.Context, summary lobby.Summary) error {
	result, err := r.serialiser.Marshall(summary)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(summary.InstanceID), result, lobbySummaryTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LobbySummaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var key = r.Key(id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LobbySummaryRepository) Key(player uuid.UUID) string {
	return strings.Join([]string{serviceKey, lobbySummaryKey, player.String()}, ":")
}
