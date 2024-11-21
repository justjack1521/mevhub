package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) InstanceID() uuid.UUID {
	return e.id
}
