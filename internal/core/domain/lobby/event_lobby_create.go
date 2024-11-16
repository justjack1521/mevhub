package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceCreatedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	quest   uuid.UUID
	party   string
	comment string
	min     int
}

func NewInstanceCreatedEvent(ctx context.Context, id, quest uuid.UUID, party, comment string, min int) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, id: id, quest: quest, party: party, comment: comment, min: min}
}

func (e InstanceCreatedEvent) Name() string {
	return "lobby.instance.created"
}

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"quest.id":   e.quest,
		"min.level":  e.min,
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceCreatedEvent) QuestID() uuid.UUID {
	return e.quest
}

func (e InstanceCreatedEvent) PartyID() string {
	return e.party
}

func (e InstanceCreatedEvent) Comment() string {
	return e.comment
}

func (e InstanceCreatedEvent) MinPlayerLevel() int {
	return e.min
}
