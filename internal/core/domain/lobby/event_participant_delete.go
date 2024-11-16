package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantDeletedEvent struct {
	ctx    context.Context
	user   uuid.UUID
	player uuid.UUID
	lobby  uuid.UUID
	slot   int
}

func NewParticipantDeletedEvent(ctx context.Context, user uuid.UUID, player uuid.UUID, lobby uuid.UUID, slot int) ParticipantDeletedEvent {
	return ParticipantDeletedEvent{ctx: ctx, user: user, player: player, lobby: lobby, slot: slot}
}

func (e ParticipantDeletedEvent) Name() string {
	return "lobby.participant.deleted"
}

func (e ParticipantDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"user.id":    e.user,
		"player.id":  e.player,
		"lobby.id":   e.lobby,
	}
}

func (e ParticipantDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantDeletedEvent) UserID() uuid.UUID {
	return e.user
}

func (e ParticipantDeletedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e ParticipantDeletedEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e ParticipantDeletedEvent) SlotIndex() int {
	return e.slot
}
