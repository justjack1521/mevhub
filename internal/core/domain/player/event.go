package player

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type DisconnectedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	user   uuid.UUID
	player uuid.UUID
	time   time.Time
}

func NewDisconnectedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID, time time.Time) DisconnectedEvent {
	return DisconnectedEvent{ctx: ctx, id: id, user: user, player: player, time: time}
}

func (e DisconnectedEvent) Name() string {
	return "player.disconnect"
}

func (e DisconnectedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":      e.Name(),
		"user.id":         e.user,
		"player.id":       e.player,
		"disconnected.at": e.time,
	}
}

func (e DisconnectedEvent) Context() context.Context {
	return e.ctx
}

func (e DisconnectedEvent) SessionID() uuid.UUID {
	return e.id
}

func (e DisconnectedEvent) UserID() uuid.UUID {
	return e.user
}

func (e DisconnectedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e DisconnectedEvent) DisconnectedAt() time.Time {
	return e.time
}

type ConnectedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	user   uuid.UUID
	player uuid.UUID
	time   time.Time
}

func NewConnectedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID, time time.Time) ConnectedEvent {
	return ConnectedEvent{ctx: ctx, id: id, user: user, player: player, time: time}
}

func (e ConnectedEvent) Name() string {
	return "player.connect"
}

func (e ConnectedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"player.id":  e.player,
		"connected.at": e.time,
	}
}

func (e ConnectedEvent) Context() context.Context {
	return e.ctx
}

func (e ConnectedEvent) SessionID() uuid.UUID {
	return e.id
}

func (e ConnectedEvent) UserID() uuid.UUID {
	return e.user
}

func (e ConnectedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e ConnectedEvent) ConnectedAt() time.Time {
	return e.time
}
