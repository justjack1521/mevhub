package session

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

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

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) UserID() uuid.UUID {
	return e.user
}

func (e InstanceCreatedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e InstanceCreatedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("user.id", e.user.String()),
		slog.String("player.id", e.player.String()),
	}
}
