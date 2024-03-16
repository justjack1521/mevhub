package session

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type InstanceReadRepository interface {
	QueryByID(ctx context.Context, id uuid.UUID) (*Instance, error)
}

type InstanceWriteRepository interface {
	Create(ctx context.Context, instance *Instance) error
	Update(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, instance *Instance) error
}

type InstanceRepository interface {
	InstanceReadRepository
	InstanceWriteRepository
}
