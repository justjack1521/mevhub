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

func NewMatchLobbyQueueRepository(client *redis.Client) *MatchLobbyQueueRepository {
	return &MatchLobbyQueueRepository{client: client}
}

func (r *MatchLobbyQueueRepository) AddActiveQuest(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID) error {
	if err := r.client.SAdd(ctx, r.activeQueueKey(mode), id.String()).Err(); err != nil {
		return err
	}
	return nil
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

	if err := r.AddActiveQuest(ctx, mode, entry.QuestID); err != nil {
		return err
	}

	var l = r.matchmakingLobbyQueueKey(mode, entry.QuestID)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	if err := r.client.ZAddArgs(ctx, l, redis.ZAddArgs{
		GT:      true,
		Members: []redis.Z{{Member: entry.LobbyID.String(), Score: float64(entry.AverageLevel)}},
	}).Err(); err != nil {
		return err
	}

	if err := r.client.ZAddArgs(ctx, t, redis.ZAddArgs{
		GT:      true,
		Members: []redis.Z{{Member: entry.LobbyID.String(), Score: float64(entry.JoinedAt.Unix())}},
	}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *MatchLobbyQueueRepository) RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
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
