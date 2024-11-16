package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantDeckChangeEvent struct {
	ctx    context.Context
	user   uuid.UUID
	player uuid.UUID
	id     uuid.UUID
	deck   int
	slot   int
}

func NewParticipantDeckChangeEvent(ctx context.Context, client, player, id uuid.UUID, deck, slot int) ParticipantDeckChangeEvent {
	return ParticipantDeckChangeEvent{ctx: ctx, user: client, player: player, id: id, deck: deck, slot: slot}
}

func (e ParticipantDeckChangeEvent) Name() string {
	return "lobby.participant.deck_change"
}

func (e ParticipantDeckChangeEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"player.id":  e.player,
		"lobby.id":   e.id,
		"deck.index": e.deck,
	}
}

func (e ParticipantDeckChangeEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantDeckChangeEvent) UserID() uuid.UUID {
	return e.user
}

func (e ParticipantDeckChangeEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e ParticipantDeckChangeEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e ParticipantDeckChangeEvent) DeckIndex() int {
	return e.deck
}

func (e ParticipantDeckChangeEvent) SlotIndex() int {
	return e.slot
}
