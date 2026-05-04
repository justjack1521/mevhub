package cache

import (
	"context"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"

	uuid "github.com/satori/go.uuid"
)

type GamePlayerLoadoutRepository struct {
	source port.GamePlayerLoadoutReadRepository
	cache  port.GamePlayerLoadoutRepository
}

func NewGamePlayerLoadoutRepository(source port.GamePlayerLoadoutReadRepository, cache port.GamePlayerLoadoutRepository) *GamePlayerLoadoutRepository {
	return &GamePlayerLoadoutRepository{source: source, cache: cache}
}

func (r *GamePlayerLoadoutRepository) Query(ctx context.Context, player uuid.UUID, index int) (game.PlayerLoadout, error) {
	hit, err := r.cache.Query(ctx, player, index)
	if err == nil {
		return hit, nil
	}

	miss, err := r.source.Query(ctx, player, index)
	if err != nil {
		return game.PlayerLoadout{}, err
	}

	_ = r.cache.Create(ctx, player, index, miss)

	return miss, nil
}

func (r *GamePlayerLoadoutRepository) Create(ctx context.Context, player uuid.UUID, index int, loadout game.PlayerLoadout) error {
	return r.cache.Create(ctx, player, index, loadout)
}

func (r *GamePlayerLoadoutRepository) Delete(ctx context.Context, player uuid.UUID, index int) error {
	return r.cache.Delete(ctx, player, index)
}
