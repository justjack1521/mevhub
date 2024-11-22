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
	user   uuid.UUID
	player uuid.UUID
}

func NewInstanceCreatedEvent(ctx context.Context, user uuid.UUID, player uuid.UUID) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, user: user, player: player}
}

func (e InstanceCreatedEvent) Name() string {
	return "session.created"
}

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"player.id":  e.player,
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) UserID() uuid.UUID {
	return e.user
}

func (e InstanceCreatedEvent) PlayerID() uuid.UUID {
	return e.player
}

type InstanceDeletedEvent struct {
	ctx    context.Context
	lobby  uuid.UUID
	gameID uuid.UUID
	user   uuid.UUID
	player uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, lobbyID uuid.UUID, gameID uuid.UUID, user uuid.UUID, player uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, lobby: lobbyID, gameID: gameID, user: user, player: player}
}

func (e InstanceDeletedEvent) Name() string {
	return "session.deleted"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.lobby,
		"user.id":    e.user,
		"player.id":  e.player,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e InstanceDeletedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e InstanceDeletedEvent) UserID() uuid.UUID {
	return e.user
}

func (e InstanceDeletedEvent) PlayerID() uuid.UUID {
	return e.player
}
