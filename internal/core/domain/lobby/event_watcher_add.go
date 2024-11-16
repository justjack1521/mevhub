package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type WatcherAddedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	user   uuid.UUID
	player uuid.UUID
}

func NewWatcherAddedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID) WatcherAddedEvent {
	return WatcherAddedEvent{ctx: ctx, id: id, user: user, player: player}
}

func (e WatcherAddedEvent) Name() string {
	return "lobby.watcher.added"
}

func (e WatcherAddedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"user.id":    e.user,
		"player.id":  e.player,
	}
}

func (e WatcherAddedEvent) Context() context.Context {
	return e.ctx
}

func (e WatcherAddedEvent) UserID() uuid.UUID {
	return e.user
}

func (e WatcherAddedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e WatcherAddedEvent) LobbyID() uuid.UUID {
	return e.id
}
