package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantReadyEvent struct {
	ctx     context.Context
	session uuid.UUID
	lobby   uuid.UUID
	deck    int
	slot    int
}

func NewParticipantReadyEvent(ctx context.Context, session uuid.UUID, lobby uuid.UUID, deck int, slot int) ParticipantReadyEvent {
	return ParticipantReadyEvent{ctx: ctx, session: session, lobby: lobby, deck: deck, slot: slot}
}

func (e ParticipantReadyEvent) Name() string {
	return "lobby.participant.ready"
}

func (e ParticipantReadyEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"session.id": e.session,
		"lobby.id":   e.lobby,
		"slot.index": e.slot,
	}
}

func (e ParticipantReadyEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantReadyEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e ParticipantReadyEvent) DeckIndex() int {
	return e.deck
}

func (e ParticipantReadyEvent) SlotIndex() int {
	return e.slot
}
