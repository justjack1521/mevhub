package player

import (
	"context"
	"log/slog"

	uuid "github.com/satori/go.uuid"
)

type DisconnectedEvent struct {
	ctx       context.Context
	id        uuid.UUID
	user      uuid.UUID
	player    uuid.UUID
	timestamp int64
}

func NewDisconnectedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID, timestamp int64) DisconnectedEvent {
	return DisconnectedEvent{ctx: ctx, id: id, user: user, player: player, timestamp: timestamp}
}

// Timestamp is the gateway emit time of the underlying ClientDisconnected, used
// to order connect/disconnect events that may arrive out of order.
func (e DisconnectedEvent) Timestamp() int64 {
	return e.timestamp
}

func (e DisconnectedEvent) Name() string {
	return "player.disconnect"
}

func (e DisconnectedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("session.id", e.id.String()),
		slog.String("user.id", e.user.String()),
		slog.String("player.id", e.player.String()),
		slog.Int64("event.timestamp", e.timestamp),
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

type ConnectedEvent struct {
	ctx       context.Context
	id        uuid.UUID
	user      uuid.UUID
	player    uuid.UUID
	timestamp int64
}

func NewConnectedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID, timestamp int64) ConnectedEvent {
	return ConnectedEvent{ctx: ctx, id: id, user: user, player: player, timestamp: timestamp}
}

// Timestamp is the gateway emit time of the underlying ClientConnected, used to
// order connect/disconnect events that may arrive out of order.
func (e ConnectedEvent) Timestamp() int64 {
	return e.timestamp
}

func (e ConnectedEvent) Name() string {
	return "player.connect"
}

func (e ConnectedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("session.id", e.id.String()),
		slog.String("user.id", e.user.String()),
		slog.String("player.id", e.player.String()),
		slog.Int64("event.timestamp", e.timestamp),
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
