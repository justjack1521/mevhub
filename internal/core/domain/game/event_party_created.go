package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type PartyCreatedEvent struct {
	ctx     context.Context
	partyID uuid.UUID
	gameID  uuid.UUID
}

func NewPartyCreatedEvent(ctx context.Context, partyID, gameID uuid.UUID) PartyCreatedEvent {
	return PartyCreatedEvent{ctx: ctx, partyID: partyID, gameID: gameID}
}

func (e PartyCreatedEvent) Name() string {
	return "game.instance.ready"
}

func (e PartyCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"party.id":   e.partyID,
		"game.id":    e.gameID,
	}
}

func (e PartyCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e PartyCreatedEvent) PartyID() uuid.UUID {
	return e.partyID
}

func (e PartyCreatedEvent) GameID() uuid.UUID {
	return e.gameID
}
