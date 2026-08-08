package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type PlayerDisconnectedEvent struct {
	ctx      context.Context
	gameID   uuid.UUID
	userID   uuid.UUID
	playerID uuid.UUID
}

func NewPlayerDisconnectedEvent(ctx context.Context, gameID, userID, playerID uuid.UUID) PlayerDisconnectedEvent {
	return PlayerDisconnectedEvent{ctx: ctx, gameID: gameID, userID: userID, playerID: playerID}
}

func (e PlayerDisconnectedEvent) Name() string {
	return "game.player.disconnected"
}

func (e PlayerDisconnectedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("game.id", e.gameID.String()),
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
	}
}

func (e PlayerDisconnectedEvent) Context() context.Context {
	return e.ctx
}

func (e PlayerDisconnectedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PlayerDisconnectedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e PlayerDisconnectedEvent) PlayerID() uuid.UUID {
	return e.playerID
}
