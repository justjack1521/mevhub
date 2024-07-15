package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/domain/game"
	"strconv"
	"strings"
	"time"
)

const (
	gameParticipantKey = "game_participant"
	gameParticipantTTL = time.Hour * 3
)

type GameParticipantRepository struct {
	client     *redis.Client
	serialiser serial.GamePlayerParticipantSerialiser
}

func NewGameParticipantRepository(client *redis.Client, serialiser serial.GamePlayerParticipantSerialiser) *GameParticipantRepository {
	return &GameParticipantRepository{client: client, serialiser: serialiser}
}

func (r *GameParticipantRepository) Query(ctx context.Context, id uuid.UUID, slot int) (game.PlayerParticipant, error) {
	value, err := r.client.Get(ctx, r.Key(id, slot)).Result()
	if err != nil {
		return game.PlayerParticipant{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return game.PlayerParticipant{}, err
	}
	return result, nil
}

func (r *GameParticipantRepository) Create(ctx context.Context, id uuid.UUID, slot int, participant game.PlayerParticipant) error {
	result, err := r.serialiser.Marshall(participant)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(id, slot), result, gameInstanceTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GameParticipantRepository) Key(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, gameInstanceKey, id.String(), strconv.Itoa(slot)}, ":")
}
