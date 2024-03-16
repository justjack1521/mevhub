package session

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type InstanceEvent interface {
	mevent.ContextEvent
	mevent.ClientEvent
	mevent.PlayerEvent
	LobbyID() uuid.UUID
}

type InstanceCreatedEvent struct {
	ctx    context.Context
	client uuid.UUID
	player uuid.UUID
}

func NewInstanceCreatedEvent(ctx context.Context, client uuid.UUID, player uuid.UUID) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, client: client, player: player}
}

func (e InstanceCreatedEvent) Name() string {
	return "session.created"
}

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"client.id":  e.client,
		"player.id":  e.player,
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) ClientID() uuid.UUID {
	return e.client
}

func (e InstanceCreatedEvent) PlayerID() uuid.UUID {
	return e.player
}

type InstanceDeletedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	client uuid.UUID
	player uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id uuid.UUID, client uuid.UUID, player uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id, client: client, player: player}
}

func (e InstanceDeletedEvent) Name() string {
	return "session.deleted"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"session.id": e.id,
		"client.id":  e.client,
		"player.id":  e.player,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) ClientID() uuid.UUID {
	return e.client
}

func (e InstanceDeletedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}
