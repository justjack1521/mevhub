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
	matchmakingLobbyPlayerActiveQuests = "matchmaking_lobby_player_active"
	matchmakingLobbyPlayerQueueKey     = "matchmaking_lobby_player_queue"
	matchmakingLobbyPlayerQueueTimeKey = "matchmaking_lobby_player_queue_time"
	matchmakingPlayerQueueKey          = "matchmaking_player_queue"
	matchmakingQueueTimeKey            = "matchmaking_player_queue_time"
)

type MatchLobbyPlayerQueueRepository struct {
	client *redis.Client
}

func NewMatchLobbyPlayerQueueRepository(client *redis.Client) *MatchLobbyPlayerQueueRepository {
	return &MatchLobbyPlayerQueueRepository{client: client}
}

func (r *MatchLobbyPlayerQueueRepository) GetCountQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) (int, error) {
	var q = r.matchmakingLobbyQueueKey(mode, quest)
	result, err := r.client.ZCard(ctx, q).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

func (r *MatchLobbyPlayerQueueRepository) GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error) {
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

func (r *MatchLobbyPlayerQueueRepository) UpdateLobbyScore(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID, score int) error {
	var l = r.matchmakingLobbyQueueKey(mode, quest)
	if err := r.client.ZAddArgs(ctx, l, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(score),
				Member: id.String(),
			},
		},
	}).Err(); err != nil {
		return err
	}
	return nil
}

func (r *MatchLobbyPlayerQueueRepository) RemoveInactiveQuest(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) error {
	if err := r.client.SRem(ctx, r.activeQueueKey(mode), quest.String()).Err(); err != nil {
		return err
	}
	return nil
}

func (r *MatchLobbyPlayerQueueRepository) AddLobbyToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry) error {

	var l = r.matchmakingLobbyQueueKey(mode, entry.QuestID)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, r.activeQueueKey(mode), entry.QuestID.String())
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

func (r *MatchLobbyPlayerQueueRepository) AddPlayerToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry) error {

	var q = r.matchmakingPlayerQueueKey(mode, entry.QuestID)
	var t = r.matchmakingPlayerQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, r.activeQueueKey(mode), entry.QuestID.String())
		pipe.ZAddArgs(ctx, q, redis.ZAddArgs{
			GT:      true,
			Members: []redis.Z{{Member: entry.UserID.String(), Score: float64(entry.DeckLevel)}},
		})
		pipe.ZAddArgs(ctx, t, redis.ZAddArgs{
			GT:      true,
			Members: []redis.Z{{Member: entry.UserID.String(), Score: float64(entry.JoinedAt.Unix())}},
		})
		return nil
	})
	return err

}

func (r *MatchLobbyPlayerQueueRepository) FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry, offset int) (match.PlayerQueueEntry, error) {

	var q = r.matchmakingPlayerQueueKey(mode, entry.QuestID)
	var t = r.matchmakingPlayerQueueTimeKey(mode)

	lower := float64(entry.AverageLevel - offset)
	upper := float64(entry.AverageLevel + offset)

	results, err := r.client.ZRangeByScoreWithScores(ctx, q, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", lower),
		Max:    fmt.Sprintf("%f", upper),
		Offset: 0,
		Count:  0,
	}).Result()

	if err != nil {
		return match.PlayerQueueEntry{}, err
	}

	var target match.PlayerQueueEntry
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
			target = match.PlayerQueueEntry{
				UserID:    id,
				QuestID:   entry.QuestID,
				DeckLevel: int(result.Score),
				JoinedAt:  time.Unix(oldest, 0),
			}
		}

	}

	return target, nil

}

func (r *MatchLobbyPlayerQueueRepository) RemoveExpiredLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) error {

	var q = r.matchmakingLobbyQueueKey(mode, quest)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	var expire = time.Now().UTC().Add(time.Minute * -20)
	var threshold = float64(expire.Unix())

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, q, "-inf", fmt.Sprintf("%f", threshold))
		pipe.ZRemRangeByScore(ctx, t, "-inf", fmt.Sprintf("%f", threshold))
		return nil
	})
	return err

}

func (r *MatchLobbyPlayerQueueRepository) GetQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) ([]match.LobbyQueueEntry, error) {

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

func (r *MatchLobbyPlayerQueueRepository) RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {

	q := r.matchmakingLobbyQueueKey(mode, quest)
	t := r.matchmakingLobbyQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, q, id.String())
		pipe.ZRem(ctx, t, id.String())
		return nil
	})
	return err

}

func (r *MatchLobbyPlayerQueueRepository) RemovePlayerFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {

	q := r.matchmakingPlayerQueueKey(mode, quest)
	t := r.matchmakingPlayerQueueTimeKey(mode)

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, q, id.String())
		pipe.ZRem(ctx, t, id.String())
		return nil
	})
	return err

}

func (r *MatchLobbyPlayerQueueRepository) activeQueueKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyPlayerActiveQuests, string(mode)}, ":")
}

func (r *MatchLobbyPlayerQueueRepository) matchmakingLobbyQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyPlayerQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchLobbyPlayerQueueRepository) matchmakingPlayerQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingPlayerQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchLobbyPlayerQueueRepository) matchmakingLobbyQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyPlayerQueueTimeKey, string(mode)}, ":")
}

func (r *MatchLobbyPlayerQueueRepository) matchmakingPlayerQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingQueueTimeKey, string(mode)}, ":")
}
