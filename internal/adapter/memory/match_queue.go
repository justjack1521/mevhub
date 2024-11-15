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
	matchmakingQueueKey     = "matchmaking_queue"
	matchmakingQueueTimeKey = "matchmaking_queue_time"
)

var (
	errFailedAddPlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed add player to matchmaking queue: %w", err)
	}
	errFailedAddPlayerTimeQueue = func(err error) error {
		return fmt.Errorf("failed add player to matchmaking time queue: %w", err)
	}
	errFailedRemovePlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed remove player to matchmaking queue: %w", err)
	}
	errFailedRemovePlayerTimeQueue = func(err error) error {
		return fmt.Errorf("failed remove player to matchmaking time queue: %w", err)
	}
	errFailedGetSizeOfMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to get size of matchmaking queue: %w", err)
	}
	errFailedRemoveActiveQuestQueue = func(err error) error {
		return fmt.Errorf("failed to remove active quest queue: %w", err)
	}
	errFailedAddActiveQuestQueue = func(err error) error {
		return fmt.Errorf("failed to add active quest queue: %w", err)
	}
	errFailedListActiveQuestQueues = func(err error) error {
		return fmt.Errorf("failed to list active quests: %w", err)
	}
	errFailedCountActiveQuestQueues = func(err error) error {
		return fmt.Errorf("failed to count active quests: %w", err)
	}
	errFailedRangePlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to range player matchmaking queue: %w", err)
	}
	errFailedListPlayerTimeQueueScores = func(err error) error {
		return fmt.Errorf("failed to list player time queue scores: %w", err)
	}
)

type MatchPlayerQueueRepository struct {
	client *redis.Client
}

func NewMatchPlayerQueueRepository(client *redis.Client) *MatchPlayerQueueRepository {
	return &MatchPlayerQueueRepository{client: client}
}

func (r *MatchPlayerQueueRepository) GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error) {
	var a = r.activeQueueKey(mode)

	ids, err := r.client.SMembers(ctx, a).Result()
	if err != nil {
		return nil, errFailedListActiveQuestQueues(err)
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

func (r *MatchPlayerQueueRepository) AddPlayerToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry) error {

	var a = r.activeQueueKey(mode)
	var q = r.matchmakingQueueKey(mode, entry.QuestID)
	var t = r.matchmakingQueueTimeKey(mode)

	if err := r.client.ZAddArgs(ctx, q, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(entry.DeckLevel),
				Member: entry.PlayerID.String(),
			},
		},
	}).Err(); err != nil {
		return errFailedAddPlayerMatchingQueue(err)
	}

	if err := r.client.ZAddArgs(ctx, t, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(time.Now().UTC().Unix()),
				Member: entry.PlayerID.String(),
			},
		},
	}).Err(); err != nil {
		return errFailedAddPlayerTimeQueue(err)
	}

	if err := r.client.SAdd(ctx, a, entry.QuestID.String()).Err(); err != nil {
		return errFailedAddActiveQuestQueue(err)
	}

	return nil

}

func (r *MatchPlayerQueueRepository) FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry, offset int) (match.PlayerQueueEntry, error) {

	var q = r.matchmakingQueueKey(mode, entry.QuestID)
	var t = r.matchmakingQueueTimeKey(mode)

	lower := float64(entry.DeckLevel - offset)
	upper := float64(entry.DeckLevel + offset)

	results, err := r.client.ZRangeByScoreWithScores(ctx, q, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", lower),
		Max:    fmt.Sprintf("%f", upper),
		Offset: 0,
		Count:  0,
	}).Result()

	if err != nil {
		return match.PlayerQueueEntry{}, errFailedRangePlayerMatchingQueue(err)
	}

	var target match.PlayerQueueEntry
	var oldest int64 = math.MaxInt64

	for _, result := range results {

		playerID, err := uuid.FromString(result.Member.(string))
		if err != nil {
			continue
		}

		joined, err := r.client.ZScore(ctx, t, playerID.String()).Result()
		if err != nil {
			continue
		}

		if int64(joined) < oldest {
			oldest = int64(joined)
			target = match.PlayerQueueEntry{
				PlayerID:  playerID,
				QuestID:   entry.QuestID,
				DeckLevel: int(result.Score),
				JoinedAt:  time.Unix(oldest, 0),
			}
		}
	}

	return target, nil

}

func (r *MatchPlayerQueueRepository) GetQueuedPlayers(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID) ([]match.PlayerQueueEntry, error) {

	var q = r.matchmakingQueueKey(mode, id)
	var t = r.matchmakingQueueTimeKey(mode)

	results, err := r.client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   q,
		Start: 0,
		Stop:  -1,
	}).Result()

	if err != nil {
		return nil, errFailedRangePlayerMatchingQueue(err)
	}

	players := make([]match.PlayerQueueEntry, 0, len(results))
	playerIDs := make([]string, len(results))

	for i, result := range results {
		playerIDs[i] = result.Member.(string)
	}

	waits, err := r.client.ZMScore(ctx, t, playerIDs...).Result()
	if err != nil {
		return nil, errFailedListPlayerTimeQueueScores(err)
	}

	for i, result := range results {
		player, err := uuid.FromString(result.Member.(string))
		if err != nil {
			continue
		}
		level := int(result.Score)
		joinedAt := time.Unix(int64(waits[i]), 0)

		entry := match.PlayerQueueEntry{
			PlayerID:  player,
			QuestID:   id,
			DeckLevel: level,
			JoinedAt:  joinedAt,
		}
		players = append(players, entry)
	}

	return players, nil

}

func (r *MatchPlayerQueueRepository) RemovePlayerFromQueue(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, player uuid.UUID) error {

	var q = r.matchmakingQueueKey(mode, id)
	var a = r.activeQueueKey(mode)

	if err := r.client.ZRem(ctx, q, player.String()).Err(); err != nil {
		return errFailedRemovePlayerMatchingQueue(err)
	}

	if err := r.client.ZRem(ctx, r.matchmakingQueueTimeKey(mode), player.String()).Err(); err != nil {
		return errFailedRemovePlayerTimeQueue(err)
	}

	size, err := r.client.ZCard(ctx, q).Result()
	if err != nil {
		return errFailedCountActiveQuestQueues(err)
	}
	if size == 0 {
		if err := r.client.SRem(ctx, a, id.String()).Err(); err != nil {
			return errFailedRemoveActiveQuestQueue(err)
		}
	}

	return nil

}

func (r *MatchPlayerQueueRepository) activeQueueKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingQueueKey, string(mode)}, ":")
}

func (r *MatchPlayerQueueRepository) matchmakingQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchPlayerQueueRepository) matchmakingQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingQueueTimeKey, string(mode)}, ":")
}
