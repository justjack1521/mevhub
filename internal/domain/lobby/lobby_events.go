package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type LobbyEvent interface {
	LobbyID() uuid.UUID
}

type InstanceCreatedEvent struct {
	ctx     context.Context
	id      uuid.UUID
	quest   uuid.UUID
	party   string
	comment string
	min     int
}

func NewInstanceCreatedEvent(ctx context.Context, id, quest uuid.UUID, party, comment string, min int) InstanceCreatedEvent {
	return InstanceCreatedEvent{ctx: ctx, id: id, quest: quest, party: party, comment: comment, min: min}
}

func (e InstanceCreatedEvent) Name() string {
	return "lobby.instance.created"
}

func (e InstanceCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"quest.id":   e.quest,
		"min.level":  e.min,
	}
}

func (e InstanceCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceCreatedEvent) LobbyID() uuid.UUID {
	return e.id
}

func (e InstanceCreatedEvent) QuestID() uuid.UUID {
	return e.quest
}

func (e InstanceCreatedEvent) PartyID() string {
	return e.party
}

func (e InstanceCreatedEvent) Comment() string {
	return e.comment
}

func (e InstanceCreatedEvent) MinPlayerLevel() int {
	return e.min
}

type InstanceDeletedEvent struct {
	ctx context.Context
	id  uuid.UUID
}

func NewInstanceDeletedEvent(ctx context.Context, id uuid.UUID) InstanceDeletedEvent {
	return InstanceDeletedEvent{ctx: ctx, id: id}
}

func (e InstanceDeletedEvent) Name() string {
	return "lobby.instance.deleted"
}

func (e InstanceDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
	}
}

func (e InstanceDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e InstanceDeletedEvent) LobbyID() uuid.UUID {
	return e.id
}

type WatcherAddedEvent struct {
	ctx    context.Context
	id     uuid.UUID
	client uuid.UUID
	player uuid.UUID
}

func NewWatcherAddedEvent(ctx context.Context, id uuid.UUID, client uuid.UUID, player uuid.UUID) WatcherAddedEvent {
	return WatcherAddedEvent{ctx: ctx, id: id, client: client, player: player}
}

func (e WatcherAddedEvent) Name() string {
	return "lobby.watcher.added"
}

func (e WatcherAddedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"lobby.id":   e.id,
		"client.id":  e.client,
		"player.id":  e.player,
	}
}

func (e WatcherAddedEvent) Context() context.Context {
	return e.ctx
}

func (e WatcherAddedEvent) ClientID() uuid.UUID {
	return e.client
}

func (e WatcherAddedEvent) PlayerID() uuid.UUID {
	return e.player
}

func (e WatcherAddedEvent) LobbyID() uuid.UUID {
	return e.id
}
