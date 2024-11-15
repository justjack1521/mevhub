package port

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
)

type MatchPlayerQueueReadRepository interface {
	GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error)
	GetQueuedPlayers(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID) ([]match.PlayerQueueEntry, error)
	FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry, offset int) (match.PlayerQueueEntry, error)
}

var (
	ErrFailedAddPlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to add player to matchmaking queue: %w", err)
	}
	ErrFailedRemovePlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to remove player from matchmaking queue: %w", err)
	}
)

type MatchPlayerQueueWriteRepository interface {
	AddPlayerToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry) error
	RemovePlayerFromQueue(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, player uuid.UUID) error
}

type MatchPlayerQueueRepository interface {
	MatchPlayerQueueReadRepository
	MatchPlayerQueueWriteRepository
}
