package port

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
)

var (
	ErrFailedAddPlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to add player to matchmaking queue: %w", err)
	}
	ErrFailedRemovePlayerMatchingQueue = func(err error) error {
		return fmt.Errorf("failed to remove player from matchmaking queue: %w", err)
	}
)

type MatchPlayerQueueReadRepository interface {
	GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error)
	GetQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID) ([]match.LobbyQueueEntry, error)
	FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry, offset int) (match.PlayerQueueEntry, error)
}

type MatchPlayerQueueWriteRepository interface {
	AddPlayerToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.PlayerQueueEntry) error
	AddLobbyToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry) error
	UpdateLobbyScore(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, score int) error
	RemovePlayerFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error
	RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error
}

type MatchPlayerQueueRepository interface {
	MatchPlayerQueueReadRepository
	MatchPlayerQueueWriteRepository
}

//type MatchLobbyQueueReadRepository interface {
//	GetActiveQuests(ctx context.Context, mode game.ModeIdentifier) ([]uuid.UUID, error)
//	GetQueuedLobbies(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID) ([]match.LobbyQueueEntry, error)
//	FindMatch(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry, offset int) (match.LobbyQueueEntry, error)
//}
//
//type MatchLobbyQueueWriteRepository interface {
//	AddLobbyToQueue(ctx context.Context, mode game.ModeIdentifier, entry match.LobbyQueueEntry) error
//	RemoveLobbyFromQueue(ctx context.Context, mode game.ModeIdentifier, quest uuid.UUID, id uuid.UUID) error
//}
//
//type MatchLobbyQueueRepository interface {
//	MatchLobbyQueueReadRepository
//	MatchLobbyQueueWriteRepository
//}
