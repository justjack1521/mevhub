package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceDeletedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	questID uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id, quest uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id, questID: quest}
}

func (e InstanceDeletedEvent) Name() string {
	return "lobby.instance.deleted"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"quest.id":   e.questID,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceDeletedEvent) QuestID() uuid.UUID {
	return e.questID
}
