package lobby

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrLobbyInstanceNil         = errors.New("game instance is nil")
	ErrLobbyInstanceNotFound    = errors.New("lobby instance not found")
	ErrFailedQueryLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to query lobby instance: %w", err)
	}
)

type InstanceReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (*Instance, error)
	QueryByPartyID(ctx context.Context, party string) (*Instance, error)
}

var (
	ErrFailedCreateLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to create lobby instance: %w", err)
	}
	ErrFailedDeleteLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to delete lobby instance: %w", err)
	}
)

type InstanceWriteRepository interface {
	Create(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type InstanceRepository interface {
	InstanceReadRepository
	InstanceWriteRepository
}
