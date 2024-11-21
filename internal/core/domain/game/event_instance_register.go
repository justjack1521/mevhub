package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceRegisteredEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceRegisteredEvent(ctx context.Context, id uuid.UUID) InstanceRegisteredEvent {
	return InstanceRegisteredEvent{ctx: ctx, id: id}
}

func (e InstanceRegisteredEvent) Name() string {
	return "game.instance.registered"
}

func (e InstanceRegisteredEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
	}
}

func (e InstanceRegisteredEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceRegisteredEvent) InstanceID() uuid.UUID {
	return e.id
}
