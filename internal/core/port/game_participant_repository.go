package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type GameParticipantReadRepository interface {
	Query(ctx context.Context, party uuid.UUID, slot int) (*game.Participant, error)
	QueryAll(ctx context.Context, party uuid.UUID) ([]*game.Participant, error)
}

type GameParticipantWriteRepository interface {
	Create(ctx context.Context, party uuid.UUID, participant *game.Participant) error
	Delete(ctx context.Context, party uuid.UUID, slot int) error
	DeleteAll(ctx context.Context, party uuid.UUID) error
}

type GameParticipantRepository interface {
	GameParticipantReadRepository
	GameParticipantWriteRepository
}
