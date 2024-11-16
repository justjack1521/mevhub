package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
	"strconv"
	"strings"
)

const (
	LobbyKeyPrefix             string = "lobby"
	LobbyKeySuffix             string = "search"
	LobbyKeyPrimarySeparator   string = "_"
	LobbyKeySecondarySeparator string = ":"
)

type LobbySearchRedisRepository struct {
	client *redis.Client
}

func NewLobbySearchRepository(client *redis.Client) *LobbySearchRedisRepository {
	return &LobbySearchRedisRepository{client: client}
}

func (r *LobbySearchRedisRepository) Query(ctx context.Context, qry lobby.SearchQuery) ([]lobby.SearchResult, error) {

	var results = make([]lobby.SearchResult, 0)

	var cache = make(map[string]bool)

	var keys = r.GenerateKeysFromQuery(qry)

	for _, key := range keys {

		result, err := r.client.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:     key,
			Start:   qry.MinimumPlayerLevel,
			Stop:    "+inf",
			ByScore: true,
			Rev:     true,
		}).Result()

		if err != nil {
			return nil, port.ErrFailedSearchForLobbies(err)
		}

		for _, value := range result {
			if cache[value] == true {
				continue
			}
			results = append(results, lobby.SearchResult{
				LobbyID: uuid.FromStringOrNil(value),
			})
		}

	}

	return results, nil

}

func (r *LobbySearchRedisRepository) Create(ctx context.Context, instance lobby.SearchEntry) error {
	var keys = r.GenerateKeysFromInstance(instance)
	_, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			if err := pipe.ZAddArgs(ctx, key, r.ZAddArgs(instance)).Err(); err != nil {
				return err
			}
			pipe.Expire(ctx, key, lobby.KeepAliveTime)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *LobbySearchRedisRepository) ZAddArgs(instance lobby.SearchEntry) redis.ZAddArgs {
	return redis.ZAddArgs{
		NX:      true,
		XX:      false,
		LT:      false,
		GT:      false,
		Ch:      false,
		Members: []redis.Z{{Score: 0, Member: instance.InstanceID.String()}},
	}
}

func (r *LobbySearchRedisRepository) GenerateKeysFromQuery(qry lobby.SearchQuery) []string {

	var identifier = r.GenerateIdentifierKey(game.ModeIdentifier(qry.ModeIdentifier))

	if len(qry.Levels) == 0 {
		return []string{identifier}
	}

	concatenated := make([]string, 0)

	if len(qry.Levels) == 0 {
		for _, category := range qry.Categories {
			concatenated = append(concatenated, category.String())
		}
	} else if len(qry.Categories) == 0 {
		for _, level := range qry.Levels {
			concatenated = append(concatenated, strconv.Itoa(level))
		}
	} else {
		for _, level := range qry.Levels {
			for _, category := range qry.Categories {
				concatenated = append(concatenated, strconv.Itoa(level)+":"+category.String())
			}
		}
	}

	var result = make([]string, len(concatenated))

	for index, value := range concatenated {
		result[index] = strings.Join([]string{identifier, value}, LobbyKeySecondarySeparator)
	}

	return result
}

func (r *LobbySearchRedisRepository) GenerateKeysFromInstance(instance lobby.SearchEntry) []string {

	var identifier = r.GenerateIdentifierKey(game.ModeIdentifier(instance.ModeIdentifier))
	var tier = strings.Join([]string{identifier, strconv.Itoa(instance.Level)}, LobbyKeySecondarySeparator)

	if len(instance.Categories) == 0 {
		return []string{tier}
	}

	var result = make([]string, len(instance.Categories))

	for index, category := range instance.Categories {
		result[index] = strings.Join([]string{tier, category.String()}, LobbyKeySecondarySeparator)
	}
	return result
}

func (r *LobbySearchRedisRepository) GenerateIdentifierKey(identifier game.ModeIdentifier) string {
	var key = strings.Join([]string{LobbyKeyPrefix, string(identifier), LobbyKeySuffix}, LobbyKeyPrimarySeparator)
	return strings.Join([]string{serviceKey, key}, LobbyKeySecondarySeparator)
}
