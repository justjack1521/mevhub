package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type ParticipantDeletedEvent struct {
	ctx     context.Context
	gameID  uuid.UUID
	partyID uuid.UUID
	userID  uuid.UUID
	slot    int
}

func NewParticipantDeletedEvent(ctx context.Context, gameID, partyID, userID uuid.UUID, slot int) ParticipantDeletedEvent {
	return ParticipantDeletedEvent{ctx: ctx, gameID: gameID, partyID: partyID, userID: userID, slot: slot}
}

func (e ParticipantDeletedEvent) Name() string {
	return "game.participant.deleted"
}

func (e ParticipantDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantDeletedEvent) GameID() uuid.UUID {
	return e.gameID
}

func (e ParticipantDeletedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e ParticipantDeletedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e ParticipantDeletedEvent) PlayerSlot() int {
	return e.slot
}
