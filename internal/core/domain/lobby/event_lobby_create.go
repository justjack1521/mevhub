package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type InstanceCreatedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	questID uuid.UUID
	party   string
	comment string
	min     int
}

func NewInstanceCreatedEvent(ctx context.Context, id, questID uuid.UUID, party, comment string, min int) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, id: id, questID: questID, party: party, comment: comment, min: min}
}

func (e InstanceCreatedEvent) Name() string {
	return "lobby.instance.created"
}

func (e InstanceCreatedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("instance.id", e.id.String()),
		slog.String("quest.id", e.questID.String()),
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceCreatedEvent) QuestID() uuid.UUID {
	return e.questID
}

func (e InstanceCreatedEvent) PartyID() string {
	return e.party
}

func (e InstanceCreatedEvent) Comment() string {
	return e.comment
}

func (e InstanceCreatedEvent) MinPlayerLevel() int {
	return e.min
}
