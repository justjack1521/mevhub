package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type InstanceReadyEvent struct {
	ctx     context.Context
	id      uuid.UUID
	questID uuid.UUID
}

func NewInstanceReadyEvent(ctx context.Context, id, questID uuid.UUID) InstanceReadyEvent {
	return InstanceReadyEvent{ctx: ctx, id: id, questID: questID}
}

func (e InstanceReadyEvent) Name() string {
	return "lobby.instance.ready"
}

func (e InstanceReadyEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("instance.id", e.id.String()),
		slog.String("quest.id", e.questID.String()),
	}
}

func (e InstanceReadyEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceReadyEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceReadyEvent) QuestID() uuid.UUID {
	return e.questID
}
