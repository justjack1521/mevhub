package memory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"math"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"strings"
	"time"
)

const (
	matchmakingLobbyActiveQuests = "matchmaking_lobby_active"
	matchmakingLobbyQueueKey     = "matchmaking_lobby_queue"
	matchmakingLobbyQueueTimeKey = "matchmaking_lobby_queue_time"
)

type MatchLobbyQueueRepository struct {
	client *redis.Client
}

func (r *MatchLobbyQueueRepository) GetCountQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) (int, error) {
	var q = r.matchmakingLobbyQueueKey(mode, quest)
	result, err := r.client.ZCard(ctx, q).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

func (r *MatchLobbyQueueRepository) RemoveExpiredLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) (int, error) {

	var q = r.matchmakingLobbyQueueKey(mode, quest)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	var expire = time.Now().UTC().Add(time.Minute * -20)
	var threshold = float64(expire.Unix())

	expired, err := r.client.ZRangeByScore(ctx, t, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%f", threshold),
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return 0, err
	}

	if len(expired) == 0 {
		return 0, nil
	}

	_, err = r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, q, expired)
		pipe.ZRem(ctx, t, expired)
		return nil
	})
	return len(expired), nil

}

func (r *MatchLobbyQueueRepository) RemoveInactiveQuest(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) error {
	if err := r.client.SRem(ctx, r.activeQueueKey(mode), quest.String()).Err(); err != nil {
		return err
	}
	return nil
}

func NewMatchLobbyQueueRepository(client *redis.Client) *MatchLobbyQueueRepository {
	return &MatchLobbyQueueRepository{client: client}
}

func (r *MatchLobbyQueueRepository) GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error) {
	var a = r.activeQueueKey(mode)

	ids, err := r.client.SMembers(ctx, a).Result()
	if err != nil {
		return nil, err
	}

	quests := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		quest, err := uuid.FromString(id)
		if err != nil {
			continue
		}
		quests = append(quests, quest)
	}

	return quests, nil
}

func (r *MatchLobbyQueueRepository) GetQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) ([]match.LobbyQueueEntry, error) {

	var q = r.matchmakingLobbyQueueKey(mode, quest)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	results, err := r.client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   q,
		Start: 0,
		Stop:  -1,
	}).Result()

	if err != nil {
		return nil, err
	}

	lobbies := make([]match.LobbyQueueEntry, 0, len(results))
	lobbyIDs := make([]string, len(results))

	if len(results) == 0 {
		return lobbies, err
	}

	for i, result := range results {
		lobbyIDs[i] = result.Member.(string)
	}

	waits, err := r.client.ZMScore(ctx, t, lobbyIDs...).Result()
	if err != nil {
		return nil, err
	}

	for i, result := range results {

		lobbyID, err := uuid.FromString(result.Member.(string))
		if err != nil {
			continue
		}
		level := int(result.Score)
		joinedAt := time.Unix(int64(waits[i]), 0)

		entry := match.LobbyQueueEntry{
			LobbyID:      lobbyID,
			QuestID:      quest,
			AverageLevel: level,
			JoinedAt:     joinedAt,
		}
		lobbies = append(lobbies, entry)
	}

	return lobbies, nil

}

func (r *MatchLobbyQueueRepository) FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry, offset int) (match.LobbyQueueEntry, error) {

	var q = r.matchmakingLobbyQueueKey(mode, entry.QuestID)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	lower := float64(entry.AverageLevel - offset)
	upper := float64(entry.AverageLevel + offset)

	results, err := r.client.ZRangeByScoreWithScores(ctx, q, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", lower),
		Max:    fmt.Sprintf("%f", upper),
		Offset: 0,
		Count:  0,
	}).Result()

	if err != nil {
		return match.LobbyQueueEntry{}, err
	}

	var target match.LobbyQueueEntry
	var oldest int64 = math.MaxInt64

	for _, result := range results {

		id, err := uuid.FromString(result.Member.(string))
		if err != nil {
			continue
		}

		if id == entry.LobbyID {
			continue
		}

		joined, err := r.client.ZScore(ctx, t, id.String()).Result()
		if err != nil {
			continue
		}

		if int64(joined) < oldest {
			oldest = int64(joined)
			target = match.LobbyQueueEntry{
				LobbyID:      id,
				QuestID:      entry.QuestID,
				AverageLevel: int(result.Score),
				JoinedAt:     time.Unix(oldest, 0),
			}
		}

	}

	return target, nil

}

func (r *MatchLobbyQueueRepository) AddLobbyToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry) error {

	var q = r.activeQueueKey(mode)
	var l = r.matchmakingLobbyQueueKey(mode, entry.QuestID)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, q, entry.QuestID.String())
		pipe.Expire(ctx, q, time.Minute*30)
		pipe.ZAddArgs(ctx, l, redis.ZAddArgs{
			GT:      true,
			Members: []redis.Z{{Member: entry.LobbyID.String(), Score: float64(entry.AverageLevel)}},
		})
		pipe.ZAddArgs(ctx, t, redis.ZAddArgs{
			GT:      true,
			Members: []redis.Z{{Member: entry.LobbyID.String(), Score: float64(entry.JoinedAt.Unix())}},
		})
		return nil
	})
	return err

}

func (r *MatchLobbyQueueRepository) RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {

	q := r.matchmakingLobbyQueueKey(mode, quest)
	t := r.matchmakingLobbyQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, q, id.String())
		pipe.ZRem(ctx, t, id.String())
		return nil
	})
	return err

}

func (r *MatchLobbyQueueRepository) activeQueueKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyActiveQuests, string(mode)}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingLobbyQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingLobbyQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyQueueTimeKey, string(mode)}, ":")
}
