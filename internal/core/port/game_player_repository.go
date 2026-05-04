package port

import (
	"context"
	"mevhub/internal/core/domain/game"

	uuid "github.com/satori/go.uuid"
)

type GamePlayerReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, slot int) (game.Player, error)
	QueryAll(ctx context.Context, id uuid.UUID) ([]game.Player, error)
}

type GamePlayerWriteRepository interface {
	Create(ctx context.Context, id uuid.UUID, slot int, participant game.Player) error
}

type GamePlayerRepository interface {
	GamePlayerReadRepository
	GamePlayerWriteRepository
}

type GamePlayerLoadoutReadRepository interface {
	Query(ctx context.Context, player uuid.UUID, index int) (game.PlayerLoadout, error)
}

type GamePlayerLoadoutWriteRepository interface {
	Create(ctx context.Context, player uuid.UUID, index int, loadout game.PlayerLoadout) error
	Delete(ctx context.Context, player uuid.UUID, index int) error
}

type GamePlayerLoadoutRepository interface {
	GamePlayerLoadoutReadRepository
	GamePlayerLoadoutWriteRepository
}
