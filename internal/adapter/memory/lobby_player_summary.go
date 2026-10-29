package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/core/domain/lobby"
	"strconv"
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

func (r *LobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID, deck int) (lobby.PlayerSummary, error) {
	value, err := r.client.Get(ctx, r.Key(id, deck)).Result()
	if err != nil {
		return lobby.PlayerSummary{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return lobby.PlayerSummary{}, err
	}
	return result, nil
}

func (r *LobbyPlayerSummaryRepository) Create(ctx context.Context, deck int, player lobby.PlayerSummary) error {
	result, err := r.serialiser.Marshall(player)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(player.Identity.PlayerID, deck), result, playerSummaryTTL).Err(); err != nil {
		return err
	}
	return nil
}

// Delete drops every deck's summary for the player. The keys are found by scanning the player's
// prefix rather than looping over deck indexes, so it holds however many decks a player has.
func (r *LobbyPlayerSummaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var keys []string
	iter := r.client.Scan(ctx, 0, r.PlayerPrefix(id)+":*", 100).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *LobbyPlayerSummaryRepository) Key(player uuid.UUID, deck int) string {
	return strings.Join([]string{r.PlayerPrefix(player), strconv.Itoa(deck)}, ":")
}

func (r *LobbyPlayerSummaryRepository) PlayerPrefix(player uuid.UUID) string {
	return strings.Join([]string{serviceKey, playerSummaryKey, player.String()}, ":")
}
