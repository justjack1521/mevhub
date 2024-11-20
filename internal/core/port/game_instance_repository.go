package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type GameInstanceReadRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*game.Instance, error)
}
type GameInstanceWriteRepository interface {
	Create(ctx context.Context, instance *game.Instance) error
}

type GameInstanceRepository interface {
	GameInstanceReadRepository
	GameInstanceWriteRepository
}

type GamePartyReadRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*game.Party, error)
	Query(ctx context.Context, game uuid.UUID, slot int) (*game.Party, error)
	QueryAll(ctx context.Context, game uuid.UUID) ([]*game.Party, error)
}

type GamePartyWriteRepository interface {
	Create(ctx context.Context, id uuid.UUID, party *game.Party) error
	Delete(ctx context.Context, id uuid.UUID, slot int) error
	DeleteAll(ctx context.Context, id uuid.UUID) error
}

type GamePartyRepository interface {
	GamePartyReadRepository
	GamePartyWriteRepository
}
