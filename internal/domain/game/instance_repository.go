package game

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrGameInstanceNotFound = func(err error) error {
		return fmt.Errorf("game instance not found: %w", err)
	}
	ErrGameInstanceNotFoundByID = func(id uuid.UUID) error {
		return fmt.Errorf("game instance not found by id: %s", id.String())
	}
)

type InstanceReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (*Instance, error)
}

var (
	ErrGameInstanceNil          = errors.New("game instance is nil")
	ErrFailedCreateGameInstance = func(err error) error {
		return fmt.Errorf("failed to create game instance: %w", err)
	}
)

type InstanceWriteRepository interface {
	Create(ctx context.Context, instance *Instance) error
}

type InstanceRepository interface {
	InstanceReadRepository
	InstanceWriteRepository
}
