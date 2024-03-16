package cache

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LazyLoadedLobbyPlayerSummaryRepository struct {
	source lobby.PlayerSummaryReadRepository
	cache  lobby.PlayerSummaryRepository
}

func NewLazyLoadedLobbyPlayerSummaryRepository(source lobby.PlayerSummaryReadRepository, cache lobby.PlayerSummaryRepository) *LazyLoadedLobbyPlayerSummaryRepository {
	return &LazyLoadedLobbyPlayerSummaryRepository{source: source, cache: cache}
}

func (r *LazyLoadedLobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID, index int) (lobby.PlayerSummary, error) {
	hit, err := r.cache.Query(ctx, id, index)
	if err != nil {
		miss, err := r.source.Query(ctx, id, index)
		if err != nil {
			return lobby.PlayerSummary{}, err
		}
		if err := r.cache.Create(ctx, miss); err != nil {
			return lobby.PlayerSummary{}, err
		}
		return miss, nil

	}
	return hit, nil
}

func (r *LazyLoadedLobbyPlayerSummaryRepository) Create(ctx context.Context, player lobby.PlayerSummary) error {
	return r.cache.Create(ctx, player)
}
