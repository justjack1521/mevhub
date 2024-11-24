package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type SummaryCreatedEvent struct {
	ctx        context.Context
	id         uuid.UUID
	questID    uuid.UUID
	level      int
	min        int
	categories []uuid.UUID
}

func NewSummaryCreatedEvent(ctx context.Context, id uuid.UUID, questID uuid.UUID, level, min int, categories []uuid.UUID) SummaryCreatedEvent {
	return SummaryCreatedEvent{ctx: ctx, id: id, questID: questID, level: level, min: min, categories: categories}
}

func (e SummaryCreatedEvent) Name() string {
	return "lobby.summary.created"
}

func (e SummaryCreatedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("instance.id", e.id.String()),
		slog.String("quest.id", e.questID.String()),
	}
}

func (e SummaryCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e SummaryCreatedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e SummaryCreatedEvent) QuestID() uuid.UUID {
	return e.questID
}

func (e SummaryCreatedEvent) Level() int {
	return e.level
}

func (e SummaryCreatedEvent) MinLevel() int {
	return e.min
}

func (e SummaryCreatedEvent) Categories() []uuid.UUID {
	return e.categories
}
