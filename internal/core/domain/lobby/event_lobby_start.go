package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceStartedEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceStartedEvent(ctx context.Context, id uuid.UUID) InstanceStartedEvent {
	return InstanceStartedEvent{ctx: ctx, id: id}
}

func (e InstanceStartedEvent) Name() string {
	return "lobby.instance.started"
}

func (e InstanceStartedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
	}
}

func (e InstanceStartedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceStartedEvent) LobbyID() uuid.UUID {
	return e.id
}
