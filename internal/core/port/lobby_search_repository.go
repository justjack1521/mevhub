package port

import (
	"context"
	"fmt"
	"mevhub/internal/core/domain/lobby"
)

var (
	ErrFailedSearchForLobbies = func(err error) error {
		return fmt.Errorf("failed to search for lobbies: %w", err)
	}
)

type LobbySearchReadRepository interface {
	Query(ctx context.Context, qry lobby.SearchQuery) ([]lobby.SearchResult, error)
}

type LobbySearchWriteRepository interface {
	Create(ctx context.Context, instance lobby.SearchEntry) error
}

type LobbySearchRepository interface {
	LobbySearchReadRepository
	LobbySearchWriteRepository
}
