package game

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceEvent interface {
	mevent.ContextEvent
	InstanceID() uuid.UUID
}

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

func (e InstanceDeletedEvent) InstanceID() uuid.UUID {
	return e.id
}
