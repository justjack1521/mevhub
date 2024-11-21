package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type PartyCreatedEvent struct {
	ctx        context.Context
	partyID    uuid.UUID
	gameID     uuid.UUID
	partyIndex int
}

func NewPartyCreatedEvent(ctx context.Context, partyID, gameID uuid.UUID, index int) PartyCreatedEvent {
	return PartyCreatedEvent{ctx: ctx, partyID: partyID, gameID: gameID, partyIndex: index}
}

func (e PartyCreatedEvent) Name() string {
	return "game.party.created"
}

func (e PartyCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e PartyCreatedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e PartyCreatedEvent) PartyIndex() int {
	return e.partyIndex
}

func (e PartyCreatedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PartyCreatedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("party.id", e.partyID.String()),
		slog.String("game.id", e.gameID.String()),
	}
}
