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

		var stop interface{} = "+inf"
		if qry.PlayerLevel > 0 {
			stop = qry.PlayerLevel
		}
		result, err := r.client.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:     key,
			Start:   0,
			Stop:    stop,
			ByScore: true,
		}).Result()

		if err != nil {
			return nil, port.ErrFailedSearchForLobbies(err)
		}

		for _, value := range result {
			// A lobby with several categories sits in several of the sets read.
			if cache[value] {
				continue
			}
			cache[value] = true
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

func (r *LobbySearchRedisRepository) Delete(ctx context.Context, instance lobby.SearchEntry) error {
	var keys = r.GenerateKeysFromInstance(instance)
	_, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			if err := pipe.ZRem(ctx, key, instance.InstanceID.String()).Err(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Prune removes an id from every set the query reads. A lobby that expires out
// of Redis emits no deletion, so its id would otherwise stay in the index for
// as long as new lobbies keep renewing the sets' TTL. Sets the query does not
// read are left for the searches that do.
func (r *LobbySearchRedisRepository) Prune(ctx context.Context, qry lobby.SearchQuery, id uuid.UUID) error {
	var keys = r.GenerateKeysFromQuery(qry)
	_, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			if err := pipe.ZRem(ctx, key, id.String()).Err(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *LobbySearchRedisRepository) ZAddArgs(instance lobby.SearchEntry) redis.ZAddArgs {
	return redis.ZAddArgs{
		NX:      true,
		XX:      false,
		LT:      false,
		GT:      false,
		Ch:      false,
		Members: []redis.Z{{Score: float64(instance.MinimumPlayerLevel), Member: instance.InstanceID.String()}},
	}
}

// Search keys. Every entry is written to one set per way it can be looked up,
// so any combination of filters reads only sets that hold all its matches:
//
//	<mode>                  every lobby in the mode       (no level, no category)
//	<mode>:<level>          every lobby in that tier      (level, no category)
//	<mode>:any:<category>   every lobby with that element (category, no level)
//	<mode>:<level>:<category>
//
// A quest with no categories is written under the "none" category, and a
// category search also reads "none": an element filter narrows quests that have
// an element, it does not hide the ones that don't.
const (
	LobbyKeyAnyLevel     string = "any"
	LobbyKeyNoneCategory string = "none"
)

func (r *LobbySearchRedisRepository) GenerateKeysFromQuery(qry lobby.SearchQuery) []string {

	var identifier = r.GenerateIdentifierKey(game.ModeIdentifier(qry.ModeIdentifier))

	var levels = make([]string, 0, len(qry.Levels))
	for _, level := range qry.Levels {
		levels = append(levels, strconv.Itoa(level))
	}

	if len(qry.Categories) == 0 {
		if len(levels) == 0 {
			return []string{identifier}
		}
		return r.joinKeys(identifier, levels)
	}

	var categories = make([]string, 0, len(qry.Categories)+1)
	for _, category := range qry.Categories {
		categories = append(categories, category.String())
	}
	categories = append(categories, LobbyKeyNoneCategory)

	if len(levels) == 0 {
		levels = []string{LobbyKeyAnyLevel}
	}

	return r.joinKeys(identifier, levels, categories)

}

func (r *LobbySearchRedisRepository) GenerateKeysFromInstance(instance lobby.SearchEntry) []string {

	var identifier = r.GenerateIdentifierKey(game.ModeIdentifier(instance.ModeIdentifier))
	var level = strconv.Itoa(instance.Level)

	var categories = make([]string, 0, len(instance.Categories))
	for _, category := range instance.Categories {
		categories = append(categories, category.String())
	}
	if len(categories) == 0 {
		categories = []string{LobbyKeyNoneCategory}
	}

	var result = []string{identifier}
	result = append(result, r.joinKeys(identifier, []string{level})...)
	result = append(result, r.joinKeys(identifier, []string{level, LobbyKeyAnyLevel}, categories)...)
	return result

}

// joinKeys returns prefix:a:b... for every combination of one value per part.
func (r *LobbySearchRedisRepository) joinKeys(prefix string, parts ...[]string) []string {
	var result = []string{prefix}
	for _, part := range parts {
		var next = make([]string, 0, len(result)*len(part))
		for _, head := range result {
			for _, value := range part {
				next = append(next, strings.Join([]string{head, value}, LobbyKeySecondarySeparator))
			}
		}
		result = next
	}
	return result
}

func (r *LobbySearchRedisRepository) GenerateIdentifierKey(identifier game.ModeIdentifier) string {
	var key = strings.Join([]string{LobbyKeyPrefix, string(identifier), LobbyKeySuffix}, LobbyKeyPrimarySeparator)
	return strings.Join([]string{serviceKey, key}, LobbyKeySecondarySeparator)
}
