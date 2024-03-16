package lobby

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ParticipantEvent interface {
	mevent.ContextEvent
	LobbyEvent
	DeckIndex() int
	SlotIndex() int
}

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

type ParticipantUnreadyEvent struct {
	ctx    context.Context
	client uuid.UUID
	lobby  uuid.UUID
	slot   int
}

func (e ParticipantUnreadyEvent) Name() string {
	return "lobby.participant.unready"
}

func (e ParticipantUnreadyEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"client.id":  e.client,
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

func NewParticipantUnreadyEvent(ctx context.Context, client uuid.UUID, lobby uuid.UUID, slot int) ParticipantUnreadyEvent {
	return ParticipantUnreadyEvent{ctx: ctx, client: client, lobby: lobby, slot: slot}
}

type ParticipantDeletedEvent struct {
	ctx    context.Context
	client uuid.UUID
	player uuid.UUID
	lobby  uuid.UUID
	slot   int
}

func NewParticipantDeletedEvent(ctx context.Context, client uuid.UUID, player uuid.UUID, lobby uuid.UUID, slot int) ParticipantDeletedEvent {
	return ParticipantDeletedEvent{ctx: ctx, client: client, player: player, lobby: lobby, slot: slot}
}

func (e ParticipantDeletedEvent) Name() string {
	return "lobby.participant.deleted"
}

func (e ParticipantDeletedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"client.id":  e.client,
		"player.id":  e.player,
		"lobby.id":   e.lobby,
	}
}

func (e ParticipantDeletedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantDeletedEvent) ClientID() uuid.UUID {
	return e.client
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

type ParticipantCreatedEvent struct {
	ctx    context.Context
	client uuid.UUID
	player uuid.UUID
	lobby  uuid.UUID
	deck   int
	slot   int
}

func NewParticipantCreatedEvent(ctx context.Context, client, player, lobby uuid.UUID, deck int, slot int) ParticipantCreatedEvent {
	return ParticipantCreatedEvent{ctx: ctx, client: client, player: player, lobby: lobby, deck: deck, slot: slot}
}

func (e ParticipantCreatedEvent) Name() string {
	return "lobby.participant.created"
}

func (e ParticipantCreatedEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"client.id":  e.client,
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

func (e ParticipantCreatedEvent) ClientID() uuid.UUID {
	return e.client
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

type ParticipantDeckChangeEvent struct {
	ctx    context.Context
	client uuid.UUID
	player uuid.UUID
	id     uuid.UUID
	deck   int
	slot   int
}

func (e ParticipantDeckChangeEvent) Name() string {
	return "lobby.participant.deck_change"
}

func (e ParticipantDeckChangeEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"client.id":  e.client,
		"player.id":  e.player,
		"lobby.id":   e.id,
		"deck.index": e.deck,
	}
}

func (e ParticipantDeckChangeEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantDeckChangeEvent) ClientID() uuid.UUID {
	return e.client
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

func NewParticipantDeckChangeEvent(ctx context.Context, client, player, id uuid.UUID, deck, slot int) ParticipantDeckChangeEvent {
	return ParticipantDeckChangeEvent{ctx: ctx, client: client, player: player, id: id, deck: deck, slot: slot}
}
