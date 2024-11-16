package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceDeletedEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id}
}

func (e InstanceDeletedEvent) Name() string {
	return "lobby.instance.deleted"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}
