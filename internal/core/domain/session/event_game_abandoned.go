package session

import (
	"context"
	"log/slog"

	uuid "github.com/satori/go.uuid"
)

// GameAbandonedEvent is a player moving on to another lobby while their session
// still seats them in a live game. The game is told to remove them, exactly as
// if their session had ended.
type GameAbandonedEvent struct {
	ctx      context.Context
	userID   uuid.UUID
	playerID uuid.UUID
	gameID   uuid.UUID
}

func NewGameAbandonedEvent(ctx context.Context, userID, playerID, gameID uuid.UUID) GameAbandonedEvent {
	return GameAbandonedEvent{ctx: ctx, userID: userID, playerID: playerID, gameID: gameID}
}

func (e GameAbandonedEvent) Name() string {
	return "session.game.abandoned"
}

func (e GameAbandonedEvent) Context() context.Context {
	return e.ctx
}

func (e GameAbandonedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e GameAbandonedEvent) PlayerID() uuid.UUID {
	return e.playerID
}

func (e GameAbandonedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e GameAbandonedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
		slog.String("game.id", e.gameID.String()),
	}
}
