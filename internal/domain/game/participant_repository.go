package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type PlayerParticipantReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, slot int) (*PlayerParticipant, error)
	QueryAll(ctx context.Context, id uuid.UUID) ([]*PlayerParticipant, error)
}

type PlayerParticipantWriteRepository interface {
	Create(ctx context.Context, id uuid.UUID, slot int, participant *PlayerParticipant) error
}

type PlayerParticipantRepository interface {
	PlayerParticipantReadRepository
	PlayerParticipantWriteRepository
}

type PlayerLoadoutReadRepository interface {
	Query(ctx context.Context, player uuid.UUID, index int) (PlayerLoadout, error)
}
