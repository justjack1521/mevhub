package lobby

import (
	"context"
	"fmt"
	"time"
)

const KeepAliveTime = time.Hour * 3

var (
	ErrFailedSearchForLobbies = func(err error) error {
		return fmt.Errorf("failed to search for lobbies: %w", err)
	}
)

type SearchReadRepository interface {
	Query(ctx context.Context, qry SearchQuery) ([]SearchResult, error)
}

type SearchWriteRepository interface {
	Create(ctx context.Context, instance SearchEntry) error
}

type SearchRepository interface {
	SearchReadRepository
	SearchWriteRepository
}
