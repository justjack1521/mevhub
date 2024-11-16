package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantCreatedEvent struct {
	ctx    context.Context
	user   uuid.UUID
	player uuid.UUID
	lobby  uuid.UUID
	deck   int
	slot   int
}

func NewParticipantCreatedEvent(ctx context.Context, user, player, lobby uuid.UUID, deck int, slot int) ParticipantCreatedEvent {
	return ParticipantCreatedEvent{ctx: ctx, user: user, player: player, lobby: lobby, deck: deck, slot: slot}
}

func (e ParticipantCreatedEvent) Name() string {
	return "lobby.participant.created"
}

func (e ParticipantCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"player.id":  e.player,
		"lobby.id":   e.lobby,
		"slot.index": e.slot,
	}
}

func (e ParticipantCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantCreatedEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e ParticipantCreatedEvent) UserID() uuid.UUID {
	return e.user
}

func (e ParticipantCreatedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e ParticipantCreatedEvent) DeckIndex() int {
	return e.deck
}

func (e ParticipantCreatedEvent) SlotIndex() int {
	return e.slot
}
