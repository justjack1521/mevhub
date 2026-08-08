package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type PlayerTimedOutEvent struct {
	ctx      context.Context
	gameID   uuid.UUID
	userID   uuid.UUID
	playerID uuid.UUID
}

func NewPlayerTimedOutEvent(ctx context.Context, gameID, userID, playerID uuid.UUID) PlayerTimedOutEvent {
	return PlayerTimedOutEvent{ctx: ctx, gameID: gameID, userID: userID, playerID: playerID}
}

func (e PlayerTimedOutEvent) Name() string {
	return "game.player.timed_out"
}

func (e PlayerTimedOutEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("game.id", e.gameID.String()),
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
	}
}

func (e PlayerTimedOutEvent) Context() context.Context {
	return e.ctx
}

func (e PlayerTimedOutEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PlayerTimedOutEvent) UserID() uuid.UUID {
	return e.userID
}

func (e PlayerTimedOutEvent) PlayerID() uuid.UUID {
	return e.playerID
}
