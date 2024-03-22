package memory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/dto"
	"mevhub/internal/domain/lobby"
	"strconv"
	"strings"
	"time"
)

const (
	playerIdentityKey = "player_summary_identity"
	playerLoadoutKey  = "player_summary_loadout"
)

var (
	ErrPlayerIdentityKeyNotFound = func(key string) error {
		return fmt.Errorf("player identity not found for key %s", key)
	}
	ErrPlayerLoadoutKeyNotFound = func(key string) error {
		return fmt.Errorf("player loadout not found for key %s", key)
	}
)

type LobbyPlayerSlotSummaryRedisRepository struct {
	client *redis.Client
}

func NewLobbyPlayerSlotSummaryRepository(client *redis.Client) *LobbyPlayerSlotSummaryRedisRepository {
	return &LobbyPlayerSlotSummaryRedisRepository{client: client}
}

func (r *LobbyPlayerSlotSummaryRedisRepository) Query(ctx context.Context, id uuid.UUID, index int) (lobby.PlayerSummary, error) {

	identity, err := r.QueryIdentity(ctx, id)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	loadout, err := r.QueryLoadout(ctx, id, index)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	return lobby.PlayerSummary{
		Identity: identity,
		Loadout:  loadout,
	}, nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) QueryIdentity(ctx context.Context, id uuid.UUID) (lobby.PlayerIdentity, error) {

	var key = r.GenerateIdentityKey(id)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return lobby.PlayerIdentity{}, err
	}

	if exists == 0 {
		return lobby.PlayerIdentity{}, ErrPlayerIdentityKeyNotFound(key)
	}

	var cmd = r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return lobby.PlayerIdentity{}, cmd.Err()
	}

	var identity = &dto.PlayerIdentityRedis{}
	if err := cmd.Scan(identity); err != nil {
		return lobby.PlayerIdentity{}, err
	}

	return identity.ToEntity(), nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) QueryLoadout(ctx context.Context, id uuid.UUID, index int) (lobby.PlayerLoadout, error) {

	var key = r.GenerateLoadoutKey(id, index)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return lobby.PlayerLoadout{}, err
	}

	if exists == 0 {
		return lobby.PlayerLoadout{}, ErrPlayerLoadoutKeyNotFound(key)
	}

	var cmd = r.client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return lobby.PlayerLoadout{}, cmd.Err()
	}

	var loadout = &dto.PlayerLoadoutRedis{}
	if err := cmd.Scan(loadout); err != nil {
		return lobby.PlayerLoadout{}, err
	}

	return loadout.ToEntity(), nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) Create(ctx context.Context, summary lobby.PlayerSummary) error {

	if err := r.CreateIdentity(ctx, summary.Identity.PlayerID, summary.Identity); err != nil {
		return err
	}

	if err := r.CreateLoadout(ctx, summary.Identity.PlayerID, summary.Loadout); err != nil {
		return err
	}

	return nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) CreateIdentity(ctx context.Context, player uuid.UUID, summary lobby.PlayerIdentity) error {

	var identity = dto.PlayerIdentityRedis{
		PlayerID:      summary.PlayerID.String(),
		PlayerName:    summary.PlayerName,
		PlayerComment: summary.PlayerComment,
		PlayerLevel:   summary.PlayerLevel,
	}

	data, err := identity.ToMapStringInterface()
	if err != nil {
		return err
	}

	if err := r.client.HSet(ctx, r.GenerateIdentityKey(player), data).Err(); err != nil {
		return err
	}

	return nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) CreateLoadout(ctx context.Context, player uuid.UUID, summary lobby.PlayerLoadout) error {

	var loadout = dto.PlayerLoadoutRedis{
		DeckIndex:       summary.DeckIndex,
		JobCardID:       summary.JobCard.JobCardID.String(),
		SubJobIndex:     summary.JobCard.SubJobIndex,
		CrownLevel:      summary.JobCard.CrownLevel,
		OverBoostLevel:  summary.JobCard.OverBoostLevel,
		WeaponID:        summary.Weapon.WeaponID.String(),
		SubWeaponUnlock: summary.Weapon.SubWeaponUnlock,
		AbilityCards:    make([]dto.PlayerAbilityCardRedis, len(summary.AbilityCards)),
	}

	for index, value := range summary.AbilityCards {
		loadout.AbilityCards[index] = dto.PlayerAbilityCardRedis{
			AbilityCardID:    value.AbilityCardID.String(),
			SlotIndex:        value.SlotIndex,
			AbilityCardLevel: value.AbilityCardLevel,
			AbilityLevel:     value.AbilityLevel,
			OverBoostLevel:   value.OverBoostLevel,
		}
	}

	data, err := loadout.ToMapStringInterface()
	if err != nil {
		return err
	}

	var key = r.GenerateLoadoutKey(player, summary.DeckIndex)

	if err := r.client.HSet(ctx, key, data).Err(); err != nil {
		return err
	}

	r.client.Expire(ctx, key, time.Minute*180)

	return nil

}

func (r *LobbyPlayerSlotSummaryRedisRepository) GenerateIdentityKey(player uuid.UUID) string {
	return strings.Join([]string{serviceKeyPrefix, playerIdentityKey, player.String()}, ":")
}

func (r *LobbyPlayerSlotSummaryRedisRepository) GenerateLoadoutKey(player uuid.UUID, index int) string {
	return strings.Join([]string{
		serviceKeyPrefix,
		playerLoadoutKey,
		player.String(),
		strconv.Itoa(index),
	}, ":")
}
