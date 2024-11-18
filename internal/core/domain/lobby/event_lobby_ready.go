package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceReadyEvent struct {
	ctx     context.Context
	id      uuid.UUID
	questID uuid.UUID
}

func NewInstanceReadyEvent(ctx context.Context, id, quest uuid.UUID) InstanceReadyEvent {
	return InstanceReadyEvent{ctx: ctx, id: id, questID: quest}
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

func (e InstanceReadyEvent) QuestID() uuid.UUID {
	return e.questID
}
