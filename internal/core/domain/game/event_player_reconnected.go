package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type PlayerReconnectedEvent struct {
	ctx      context.Context
	gameID   uuid.UUID
	userID   uuid.UUID
	playerID uuid.UUID
}

func NewPlayerReconnectedEvent(ctx context.Context, gameID, userID, playerID uuid.UUID) PlayerReconnectedEvent {
	return PlayerReconnectedEvent{ctx: ctx, gameID: gameID, userID: userID, playerID: playerID}
}

func (e PlayerReconnectedEvent) Name() string {
	return "game.player.reconnected"
}

func (e PlayerReconnectedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("game.id", e.gameID.String()),
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
	}
}

func (e PlayerReconnectedEvent) Context() context.Context {
	return e.ctx
}

func (e PlayerReconnectedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PlayerReconnectedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e PlayerReconnectedEvent) PlayerID() uuid.UUID {
	return e.playerID
}
