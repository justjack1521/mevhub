package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/match"
)

type MatchInstanceReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (*match.Instance, error)
}

type MatchInstanceWriteRepository interface {
	Create(ctx context.Context, instance *match.Instance) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MatchInstanceRepository interface {
	MatchInstanceReadRepository
	MatchInstanceWriteRepository
}
