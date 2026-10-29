package port

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
)

type LobbySummaryReadRepository interface {
	Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error)
}

type LobbySearchSummaryRepository interface {
	Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error)
	QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error)
}

var (
	ErrFailedCreateLobbySummary = func(summary lobby.Summary, err error) error {
		return fmt.Errorf("failed to create summary for lobby: %s: %w", summary.InstanceID, err)
	}
	ErrFailedDeleteLobbySummary = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to delete summary for lobby: %s: %w", id, err)
	}
)

type LobbySummaryWriteRepository interface {
	Create(ctx context.Context, summary lobby.Summary) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LobbySummaryRepository interface {
	LobbySummaryReadRepository
	LobbySummaryWriteRepository
}

var (
	ErrLobbyPlayerSummaryNotFound = func(id uuid.UUID) error {
		return fmt.Errorf("lobby player %s summary not found", id.String())
	}
)

// LobbyPlayerSummaryReadRepository is keyed by deck as well as player: the summary carries the
// loadout, and a lobby shows the deck each player joined with, not whichever deck the source
// would pick by default.
type LobbyPlayerSummaryReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, deck int) (lobby.PlayerSummary, error)
}

var (
	ErrFailedCreatePlayerSummary = func(id uuid.UUID) error {
		return fmt.Errorf("failed create lobby player %s summary", id.String())
	}
)

type LobbyPlayerSummaryWriteRepository interface {
	Create(ctx context.Context, deck int, player lobby.PlayerSummary) error
	// Delete drops the player's summary for every deck.
	Delete(ctx context.Context, id uuid.UUID) error
}

type LobbyPlayerSummaryRepository interface {
	LobbyPlayerSummaryReadRepository
	LobbyPlayerSummaryWriteRepository
}
