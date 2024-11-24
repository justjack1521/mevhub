package session

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type InstanceDeletedEvent struct {
	ctx      context.Context
	userID   uuid.UUID
	playerID uuid.UUID
	lobbyID  uuid.UUID
	gameID   uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, userID, playerID, lobbyID, gameID uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, userID: userID, playerID: playerID, lobbyID: lobbyID, gameID: gameID}
}

func (e InstanceDeletedEvent) Name() string {
	return "session.deleted"
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e InstanceDeletedEvent) PlayerID() uuid.UUID {
	return e.playerID
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.lobbyID
}

func (e InstanceDeletedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e InstanceDeletedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
		slog.String("lobby.id", e.lobbyID.String()),
		slog.String("game.id", e.gameID.String()),
	}
}
