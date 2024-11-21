package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceReadyEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceReadyEvent(ctx context.Context, id uuid.UUID) InstanceReadyEvent {
	return InstanceReadyEvent{ctx: ctx, id: id}
}

func (e InstanceReadyEvent) Name() string {
	return "game.instance.ready"
}

func (e InstanceReadyEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
	}
}

func (e InstanceReadyEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceReadyEvent) InstanceID() uuid.UUID {
	return e.id
}
