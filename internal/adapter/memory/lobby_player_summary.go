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
	playerSummaryKey = "lobby_player_summary"
	playerSummaryTTL = time.Minute * 180
)

type LobbyPlayerSummaryRepository struct {
	client     *redis.Client
	serialiser serial.LobbyPlayerSummarySerialiser
}

func NewLobbyPlayerSummaryRepository(client *redis.Client, serialiser serial.LobbyPlayerSummarySerialiser) *LobbyPlayerSummaryRepository {
	return &LobbyPlayerSummaryRepository{client: client, serialiser: serialiser}
}

func (r *LobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID) (lobby.PlayerSummary, error) {
	value, err := r.client.Get(ctx, r.Key(id)).Result()
	if err != nil {
		return lobby.PlayerSummary{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return lobby.PlayerSummary{}, err
	}
	return result, nil
}

func (r *LobbyPlayerSummaryRepository) Create(ctx context.Context, player lobby.PlayerSummary) error {
	result, err := r.serialiser.Marshall(player)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(player.Identity.PlayerID), result, playerSummaryTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LobbyPlayerSummaryRepository) Key(player uuid.UUID) string {
	return strings.Join([]string{serviceKey, playerSummaryKey, player.String()}, ":")
}
