package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
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
