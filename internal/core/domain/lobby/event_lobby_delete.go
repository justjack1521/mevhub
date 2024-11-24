package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type InstanceDeletedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	questID uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id, quest uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id, questID: quest}
}

func (e InstanceDeletedEvent) Name() string {
	return "lobby.instance.deleted"
}

func (e InstanceDeletedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("instance.id", e.id.String()),
		slog.String("quest.id", e.questID.String()),
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceDeletedEvent) QuestID() uuid.UUID {
	return e.questID
}
