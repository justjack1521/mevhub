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
	return r.query(ctx, r.Key(id, slot))
}

func (r *GameParticipantRepository) QueryAll(ctx context.Context, id uuid.UUID) ([]game.PlayerParticipant, error) {
	keys, err := r.client.Keys(ctx, r.GameKey(id)).Result()
	if err != nil {
		return nil, err
	}
	var participants = make([]game.PlayerParticipant, len(keys))
	for index, key := range keys {
		participant, err := r.query(ctx, key)
		if err != nil {
			return nil, err
		}
		participants[index] = participant
	}
	return participants, nil
}

func (r *GameParticipantRepository) Create(ctx context.Context, id uuid.UUID, slot int, participant game.PlayerParticipant) error {
	result, err := r.serialiser.Marshall(participant)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, r.Key(id, slot), result, gameParticipantTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GameParticipantRepository) query(ctx context.Context, key string) (game.PlayerParticipant, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return game.PlayerParticipant{}, err
	}
	result, err := r.serialiser.Unmarshall([]byte(value))
	if err != nil {
		return game.PlayerParticipant{}, err
	}
	return result, nil
}

func (r *GameParticipantRepository) GameKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, gameParticipantKey, id.String(), "*"}, ":")
}

func (r *GameParticipantRepository) Key(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, gameParticipantKey, id.String(), strconv.Itoa(slot)}, ":")
}
