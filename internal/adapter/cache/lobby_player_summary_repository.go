package cache

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyPlayerSummaryRepository struct {
	source port.LobbyPlayerSummaryReadRepository
	cache  port.LobbyPlayerSummaryRepository
}

func NewLobbyPlayerSummaryRepository(source port.LobbyPlayerSummaryReadRepository, cache port.LobbyPlayerSummaryRepository) *LobbyPlayerSummaryRepository {
	return &LobbyPlayerSummaryRepository{source: source, cache: cache}
}

func (r *LobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID, deck int) (lobby.PlayerSummary, error) {
	hit, err := r.cache.Query(ctx, id, deck)
	if err == nil {
		return hit, nil
	}

	miss, err := r.source.Query(ctx, id, deck)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	_ = r.cache.Create(ctx, deck, miss)

	return miss, nil
}

func (r *LobbyPlayerSummaryRepository) Create(ctx context.Context, deck int, player lobby.PlayerSummary) error {
	return r.cache.Create(ctx, deck, player)
}

func (r *LobbyPlayerSummaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cache.Delete(ctx, id)
}
