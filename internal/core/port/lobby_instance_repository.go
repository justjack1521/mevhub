package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
)

type LobbyInstanceReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (*lobby.Instance, error)
	QueryByPartyID(ctx context.Context, party string) (*lobby.Instance, error)
}

type LobbyInstanceWriteRepository interface {
	Create(ctx context.Context, instance *lobby.Instance) error
	Delete(ctx context.Context, instance *lobby.Instance) error
}

type LobbyInstanceRepository interface {
	LobbyInstanceReadRepository
	LobbyInstanceWriteRepository
}
