package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type SummaryReadRepository interface {
	QueryByID(id uuid.UUID) (Summary, error)
}

type SummaryWriteRepository interface {
	Create(ctx context.Context, summary Summary) error
}

type SummaryRepository interface {
	SummaryReadRepository
	SummaryWriteRepository
}

type PlayerSummaryReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, index int) (PlayerSummary, error)
}

type PlayerSummaryWriteRepository interface {
}

type PlayerSummaryRepository interface {
	PlayerSummaryReadRepository
	PlayerSummaryWriteRepository
}
