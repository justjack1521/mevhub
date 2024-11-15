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

type ParticipantCreatedEvent struct {
	ctx  context.Context
	id   uuid.UUID
	slot int
}

func NewParticipantCreatedEvent(ctx context.Context, id uuid.UUID, slot int) ParticipantCreatedEvent {
	return ParticipantCreatedEvent{ctx: ctx, id: id, slot: slot}
}

func (e ParticipantCreatedEvent) Name() string {
	return "game.participant.created"
}

func (e ParticipantCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":  e.Name(),
		"instance.id": e.id,
		"player.slot": e.slot,
	}
}

func (e ParticipantCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantCreatedEvent) InstanceID() uuid.UUID {
	return e.id
}

func (e ParticipantCreatedEvent) PlayerSlot() int {
	return e.slot
}
