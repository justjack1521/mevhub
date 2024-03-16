package game

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceEvent interface {
	mevent.ContextEvent
	LobbyID() uuid.UUID
}

type InstanceCreatedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	quest   uuid.UUID
	party   string
	host    uuid.UUID
	comment string
}

func NewInstanceCreatedEvent(ctx context.Context, id, quest uuid.UUID, party, comment string) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, id: id, quest: quest, party: party, comment: comment}
}

func (e InstanceCreatedEvent) Name() string {
	return "game.instance.created"
}

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
		"quest.id":    e.quest,
		"party.id":    e.party,
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

type InstanceDeletedEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id}
}

func (e InstanceDeletedEvent) Name() string {
	return "game.instance.delete"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}
