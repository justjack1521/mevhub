package lobby

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type SummaryReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (Summary, error)
	QueryByPartyID(ctx context.Context, party string) (Summary, error)
}

type SummaryWriteRepository interface {
	Create(ctx context.Context, summary Summary) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SummaryRepository interface {
	SummaryReadRepository
	SummaryWriteRepository
}

var (
	ErrLobbyPlayerSummaryNotFound = func(id uuid.UUID) error {
		return fmt.Errorf("lobby player %s summary not found", id.String())
	}
)

type PlayerSummaryReadRepository interface {
	Query(ctx context.Context, id uuid.UUID, index int) (PlayerSummary, error)
}

var (
	ErrFailedCreatePlayerSummary = func(id uuid.UUID) error {
		return fmt.Errorf("failed create lobby player %s summary", id.String())
	}
)

type PlayerSummaryWriteRepository interface {
	Create(ctx context.Context, player PlayerSummary) error
}

type PlayerSummaryRepository interface {
	PlayerSummaryReadRepository
	PlayerSummaryWriteRepository
}
