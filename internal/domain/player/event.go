package player

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type DisconnectedEvent struct {
	ctx    context.Context
	user   uuid.UUID
	player uuid.UUID
	time   time.Time
}

func NewDisconnectedEvent(ctx context.Context, user uuid.UUID, player uuid.UUID, time time.Time) DisconnectedEvent {
	return DisconnectedEvent{ctx: ctx, user: user, player: player, time: time}
}

func (e DisconnectedEvent) Name() string {
	return "player.disconnect"
}

func (e DisconnectedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name":      e.Name(),
		"user.id":         e.user,
		"player.id":       e.player,
		"disconnected.at": e.time,
	}
}

func (e DisconnectedEvent) Context() context.Context {
	return e.ctx
}

func (e DisconnectedEvent) UserID() uuid.UUID {
	return e.user
}

func (e DisconnectedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e DisconnectedEvent) DisconnectedAt() time.Time {
	return e.time
}
