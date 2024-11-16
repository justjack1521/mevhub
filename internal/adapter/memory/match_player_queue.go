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
	matchmakingLobbyQueueKey     = "matchmaking_lobby_queue"
	matchmakingLobbyQueueTimeKey = "matchmaking_lobby_queue_time"
	matchmakingPlayerQueueKey    = "matchmaking_player_queue"
	matchmakingQueueTimeKey      = "matchmaking_player_queue_time"
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

type MatchLobbyQueueRepository struct {
	client *redis.Client
}

func NewMatchPlayerQueueRepository(client *redis.Client) *MatchLobbyQueueRepository {
	return &MatchLobbyQueueRepository{client: client}
}

func (r *MatchLobbyQueueRepository) GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error) {
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

func (r *MatchLobbyQueueRepository) UpdateLobbyScore(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID, score int) error {
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
		return errFailedAddPlayerMatchingQueue(err)
	}
	return nil
}

func (r *MatchLobbyQueueRepository) AddLobbyToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry) error {

	var a = r.activeQueueKey(mode)
	var l = r.matchmakingLobbyQueueKey(mode, entry.QuestID)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	if err := r.client.ZAddArgs(ctx, l, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(entry.AverageLevel),
				Member: entry.LobbyID.String(),
			},
		},
	}).Err(); err != nil {
		return errFailedAddPlayerMatchingQueue(err)
	}

	if err := r.client.ZAddArgs(ctx, t, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(entry.JoinedAt.Unix()),
				Member: entry.LobbyID.String(),
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

func (r *MatchLobbyQueueRepository) AddPlayerToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry) error {

	var a = r.activeQueueKey(mode)
	var q = r.matchmakingPlayerQueueKey(mode, entry.QuestID)
	var t = r.matchmakingPlayerQueueTimeKey(mode)

	if err := r.client.ZAddArgs(ctx, q, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(entry.DeckLevel),
				Member: entry.UserID.String(),
			},
		},
	}).Err(); err != nil {
		return errFailedAddPlayerMatchingQueue(err)
	}

	if err := r.client.ZAddArgs(ctx, t, redis.ZAddArgs{
		GT: true,
		Members: []redis.Z{
			{
				Score:  float64(entry.JoinedAt.Unix()),
				Member: entry.UserID.String(),
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

func (r *MatchLobbyQueueRepository) FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry, offset int) (match.PlayerQueueEntry, error) {

	var p = r.matchmakingPlayerQueueKey(mode, entry.QuestID)
	var t = r.matchmakingPlayerQueueTimeKey(mode)

	lower := float64(entry.AverageLevel - offset)
	upper := float64(entry.AverageLevel + offset)

	results, err := r.client.ZRangeByScoreWithScores(ctx, p, &redis.ZRangeBy{
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

		userID, err := uuid.FromString(result.Member.(string))
		if err != nil {
			continue
		}

		joined, err := r.client.ZScore(ctx, t, userID.String()).Result()
		if err != nil {
			continue
		}

		if int64(joined) < oldest {
			oldest = int64(joined)
			target = match.PlayerQueueEntry{
				UserID:    userID,
				QuestID:   entry.QuestID,
				DeckLevel: int(result.Score),
				JoinedAt:  time.Unix(oldest, 0),
			}
		}

	}

	return target, nil

}

func (r *MatchLobbyQueueRepository) GetQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID) ([]match.LobbyQueueEntry, error) {

	var l = r.matchmakingLobbyQueueKey(mode, id)
	var t = r.matchmakingLobbyQueueTimeKey(mode)

	results, err := r.client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   l,
		Start: 0,
		Stop:  -1,
	}).Result()

	if err != nil {
		return nil, errFailedRangePlayerMatchingQueue(err)
	}

	lobbies := make([]match.LobbyQueueEntry, 0, len(results))
	lobbyIDs := make([]string, len(results))

	for i, result := range results {
		lobbyIDs[i] = result.Member.(string)
	}

	waits, err := r.client.ZMScore(ctx, t, lobbyIDs...).Result()
	if err != nil {
		return nil, errFailedListPlayerTimeQueueScores(err)
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
			QuestID:      id,
			AverageLevel: level,
			JoinedAt:     joinedAt,
		}
		lobbies = append(lobbies, entry)
	}

	return lobbies, nil

}

func (r *MatchLobbyQueueRepository) RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {

	if err := r.client.ZRem(ctx, r.matchmakingLobbyQueueKey(mode, quest), id.String()).Err(); err != nil {
		return errFailedRemovePlayerMatchingQueue(err)
	}

	if err := r.client.ZRem(ctx, r.matchmakingLobbyQueueTimeKey(mode), id.String()).Err(); err != nil {
		return errFailedRemovePlayerTimeQueue(err)
	}

	return nil

}

func (r *MatchLobbyQueueRepository) RemovePlayerFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error {

	if err := r.client.ZRem(ctx, r.matchmakingPlayerQueueKey(mode, quest), id.String()).Err(); err != nil {
		return errFailedRemovePlayerMatchingQueue(err)
	}

	if err := r.client.ZRem(ctx, r.matchmakingPlayerQueueTimeKey(mode), id.String()).Err(); err != nil {
		return errFailedRemovePlayerTimeQueue(err)
	}

	return nil

}

func (r *MatchLobbyQueueRepository) activeQueueKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingPlayerQueueKey, string(mode)}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingLobbyQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingPlayerQueueKey(mode game.ModeIdentifier, id uuid.UUID) string {
	return strings.Join([]string{serviceKey, matchmakingPlayerQueueKey, string(mode), id.String()}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingLobbyQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingLobbyQueueTimeKey, string(mode)}, ":")
}

func (r *MatchLobbyQueueRepository) matchmakingPlayerQueueTimeKey(mode game.ModeIdentifier) string {
	return strings.Join([]string{serviceKey, matchmakingQueueTimeKey, string(mode)}, ":")
}
