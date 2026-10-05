package port

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
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

// LobbySearchIndexRepository is the search side's view of the index: query it,
// and drop ids it turns out to hold for lobbies that no longer exist.
type LobbySearchIndexRepository interface {
	LobbySearchReadRepository
	Prune(ctx context.Context, qry lobby.SearchQuery, id uuid.UUID) error
}

type LobbySearchWriteRepository interface {
	Create(ctx context.Context, instance lobby.SearchEntry) error
	Delete(ctx context.Context, instance lobby.SearchEntry) error
}

type LobbySearchRepository interface {
	LobbySearchIndexRepository
	LobbySearchWriteRepository
}
