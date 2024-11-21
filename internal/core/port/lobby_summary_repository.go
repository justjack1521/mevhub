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

type LobbyPlayerSummaryReadRepository interface {
	Query(ctx context.Context, id uuid.UUID) (lobby.PlayerSummary, error)
}

var (
	ErrFailedCreatePlayerSummary = func(id uuid.UUID) error {
		return fmt.Errorf("failed create lobby player %s summary", id.String())
	}
)

type LobbyPlayerSummaryWriteRepository interface {
	Create(ctx context.Context, player lobby.PlayerSummary) error
}

type LobbyPlayerSummaryRepository interface {
	LobbyPlayerSummaryReadRepository
	LobbyPlayerSummaryWriteRepository
}
