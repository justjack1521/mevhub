package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantUnreadyEvent struct {
	ctx   context.Context
	user  uuid.UUID
	lobby uuid.UUID
	slot  int
}

func NewParticipantUnreadyEvent(ctx context.Context, user uuid.UUID, lobby uuid.UUID, slot int) ParticipantUnreadyEvent {
	return ParticipantUnreadyEvent{ctx: ctx, user: user, lobby: lobby, slot: slot}
}

func (e ParticipantUnreadyEvent) Name() string {
	return "lobby.participant.unready"
}

func (e ParticipantUnreadyEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"lobby.id":   e.lobby,
		"slot.index": e.slot,
	}
}

func (e ParticipantUnreadyEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantUnreadyEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e ParticipantUnreadyEvent) SlotIndex() int {
	return e.slot
}
