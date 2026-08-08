package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/session"
)

type SessionInstanceReadRepository interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	QueryByID(ctx context.Context, id uuid.UUID) (*session.Instance, error)
}

type SessionInstanceWriteRepository interface {
	Create(ctx context.Context, instance *session.Instance) error
	Update(ctx context.Context, instance *session.Instance) error
	// UpdateConnectionState persists only the connect/disconnect bookkeeping
	// fields (DisconnectSessionID, LastConnEventAt, DisconnectedAt) so the
	// event handlers cannot clobber concurrent lobby/game field writes with a
	// stale full-record read.
	UpdateConnectionState(ctx context.Context, instance *session.Instance) error
	Delete(ctx context.Context, instance *session.Instance) error
}

type SessionInstanceRepository interface {
	SessionInstanceReadRepository
	SessionInstanceWriteRepository
}
