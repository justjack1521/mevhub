package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type InstanceStartedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	gameID uuid.UUID
}

func NewInstanceStartedEvent(ctx context.Context, id uuid.UUID, gameID uuid.UUID) InstanceStartedEvent {
	return InstanceStartedEvent{ctx: ctx, id: id, gameID: gameID}
}

func (e InstanceStartedEvent) Name() string {
	return "lobby.instance.started"
}

func (e InstanceStartedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("instance.id", e.id.String()),
		slog.String("game.id", e.gameID.String()),
	}
}

func (e InstanceStartedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceStartedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceStartedEvent) GameID() uuid.UUID {
	return e.gameID
}
