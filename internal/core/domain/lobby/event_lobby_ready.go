package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceReadyEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func (e InstanceReadyEvent) Name() string {
	return "lobby.instance.ready"
}

func (e InstanceReadyEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
	}
}

func (e InstanceReadyEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceReadyEvent) LobbyID() uuid.UUID {
	return e.id
}
