package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type InstanceCreatedEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceCreatedEvent(ctx context.Context, id uuid.UUID) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, id: id}
}

func (e InstanceCreatedEvent) Name() string {
	return "game.instance.created"
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) InstanceID() uuid.UUID {
	return e.id
}
