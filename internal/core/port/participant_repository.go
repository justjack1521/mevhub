package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type PlayerParticipantReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, slot int) (*game.PlayerParticipant, error)
	QueryAll(ctx context.Context, id uuid.UUID) ([]*game.PlayerParticipant, error)
}

type PlayerParticipantWriteRepository interface {
	Create(ctx context.Context, id uuid.UUID, slot int, participant *game.PlayerParticipant) error
}

type PlayerParticipantRepository interface {
	PlayerParticipantReadRepository
	PlayerParticipantWriteRepository
}

type PlayerLoadoutReadRepository interface {
	Query(ctx context.Context, player uuid.UUID, index int) (game.PlayerLoadout, error)
}
