package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type SummaryCreatedEvent struct {
	ctx        context.Context
	id         uuid.UUID
	mode       string
	level      int
	min        int
	categories []uuid.UUID
}

func NewSummaryCreatedEvent(ctx context.Context, id uuid.UUID, mode string, level, min int, categories []uuid.UUID) SummaryCreatedEvent {
	return SummaryCreatedEvent{ctx: ctx, id: id, mode: mode, level: level, min: min, categories: categories}
}

func (e SummaryCreatedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e SummaryCreatedEvent) Name() string {
	return "lobby.summary.created"
}

func (e SummaryCreatedEvent) Mode() string {
	return e.mode
}

func (e SummaryCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"mode":       e.mode,
	}
}

func (e SummaryCreatedEvent) Context() context.Context {
	return e.ctx
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
