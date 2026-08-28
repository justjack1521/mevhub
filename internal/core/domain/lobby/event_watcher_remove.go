package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type WatcherRemovedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	user   uuid.UUID
	player uuid.UUID
}

func NewWatcherRemovedEvent(ctx context.Context, id uuid.UUID, user uuid.UUID, player uuid.UUID) WatcherRemovedEvent {
	return WatcherRemovedEvent{ctx: ctx, id: id, user: user, player: player}
}

func (e WatcherRemovedEvent) Name() string {
	return "lobby.watcher.removed"
}

func (e WatcherRemovedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"user.id":    e.user,
		"player.id":  e.player,
	}
}

func (e WatcherRemovedEvent) Context() context.Context {
	return e.ctx
}

func (e WatcherRemovedEvent) UserID() uuid.UUID {
	return e.user
}

func (e WatcherRemovedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e WatcherRemovedEvent) LobbyID() uuid.UUID {
	return e.id
}
