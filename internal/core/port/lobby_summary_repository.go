package port

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
)

type LobbySummaryReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (lobby.Summary, error)
	QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error)
}

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
