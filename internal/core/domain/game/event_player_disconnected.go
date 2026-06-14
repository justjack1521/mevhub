package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type PlayerDisconnectedEvent struct {
	ctx      context.Context
	gameID   uuid.UUID
	partyID  uuid.UUID
	userID   uuid.UUID
	playerID uuid.UUID
}

func NewPlayerDisconnectedEvent(ctx context.Context, gameID, partyID, userID, playerID uuid.UUID) PlayerDisconnectedEvent {
	return PlayerDisconnectedEvent{ctx: ctx, gameID: gameID, partyID: partyID, userID: userID, playerID: playerID}
}

func (e PlayerDisconnectedEvent) Name() string {
	return "game.player.disconnected"
}

func (e PlayerDisconnectedEvent) Context() context.Context {
	return e.ctx
}

func (e PlayerDisconnectedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e PlayerDisconnectedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e PlayerDisconnectedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e PlayerDisconnectedEvent) PlayerID() uuid.UUID {
	return e.playerID
}
