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
	gameParticipantKey          = "game_participant"
	gameParticipantKeySeparator = ":"
	gameParticipantTTL          = time.Hour * 3
)

type GameParticipantRepository struct {
	client *redis.Client
}

func NewGameParticipantRepository(client *redis.Client) *GameParticipantRepository {
	return &GameParticipantRepository{client: client}
}

func (r *GameParticipantRepository) query(ctx context.Context, key string) (*game.Participant, error) {
	var cmd = r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var result = &dto.GameParticipantRedis{}
	if err := cmd.Scan(result); err != nil {
		return nil, err
	}
	return result.ToEntity(), nil
}

func (r *GameParticipantRepository) Query(ctx context.Context, party uuid.UUID, slot int) (*game.Participant, error) {
	var key = r.Key(party, slot)
	return r.query(ctx, key)

}

func (r *GameParticipantRepository) QueryAll(ctx context.Context, party uuid.UUID) ([]*game.Participant, error) {
	keys, err := r.client.Keys(ctx, r.PartyKey(party)).Result()
	if err != nil {
		return nil, err
	}
	var participants = make([]*game.Participant, len(keys))
	for index, key := range keys {
		participant, err := r.query(ctx, key)
		if err != nil {
			return nil, err
		}
		participants[index] = participant
	}
	return participants, nil
}

func (r *GameParticipantRepository) Create(ctx context.Context, party uuid.UUID, participant *game.Participant) error {

	var result = &dto.GameParticipantRedis{
		UserID:     participant.UserID.String(),
		PlayerID:   participant.PlayerID.String(),
		PlayerSlot: participant.PlayerSlot,
		DeckIndex:  participant.DeckIndex,
		BotControl: participant.BotControl,
	}

	var key = r.Key(party, participant.PlayerSlot)

	if err := r.client.HSet(ctx, key, result.ToMapStringInterface()).Err(); err != nil {
		return err
	}

	r.client.Expire(ctx, key, gameParticipantTTL)
	return nil

}

func (r *GameParticipantRepository) Delete(ctx context.Context, party uuid.UUID, slot int) error {
	var key = r.Key(party, slot)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GameParticipantRepository) DeleteAll(ctx context.Context, party uuid.UUID) error {

	keys, err := r.client.Keys(ctx, r.PartyKey(party)).Result()
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

func (r *GameParticipantRepository) PartyKey(id uuid.UUID) string {
	return strings.Join([]string{serviceKey, gameParticipantKey, id.String(), "*"}, gameParticipantKeySeparator)
}

func (r *GameParticipantRepository) Key(id uuid.UUID, slot int) string {
	return strings.Join([]string{serviceKey, gameParticipantKey, id.String(), strconv.Itoa(slot)}, gameParticipantKeySeparator)
}
