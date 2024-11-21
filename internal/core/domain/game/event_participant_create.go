package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type ParticipantCreatedEvent struct {
	ctx     context.Context
	gameID  uuid.UUID
	partyID uuid.UUID
	userID  uuid.UUID
	slot    int
}

func NewParticipantCreatedEvent(ctx context.Context, gameID, partyID, userID uuid.UUID, slot int) ParticipantCreatedEvent {
	return ParticipantCreatedEvent{ctx: ctx, gameID: gameID, partyID: partyID, userID: userID, slot: slot}
}

func (e ParticipantCreatedEvent) Name() string {
	return "game.participant.created"
}

func (e ParticipantCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantCreatedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e ParticipantCreatedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e ParticipantCreatedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e ParticipantCreatedEvent) PlayerSlot() int {
	return e.slot
}
