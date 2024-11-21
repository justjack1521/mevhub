package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type PartyDeletedEvent struct {
	ctx     context.Context
	partyID uuid.UUID
	gameID  uuid.UUID
}

func NewPartyDeletedEvent(ctx context.Context, partyID uuid.UUID, gameID uuid.UUID) PartyDeletedEvent {
	return PartyDeletedEvent{ctx: ctx, partyID: partyID, gameID: gameID}
}

func (e PartyDeletedEvent) Name() string {
	return "game.party.deleted"
}

func (e PartyDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e PartyDeletedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e PartyDeletedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PartyDeletedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("party.id", e.partyID.String()),
		slog.String("game.id", e.gameID.String()),
	}
}
