package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type ParticipantCreatedEvent struct {
	ctx      context.Context
	userID   uuid.UUID
	playerID uuid.UUID
	lobbyID  uuid.UUID
	deck     int
	slot     int
}

func NewParticipantCreatedEvent(ctx context.Context, userID, playerID, lobbyID uuid.UUID, deck int, slot int) ParticipantCreatedEvent {
	return ParticipantCreatedEvent{ctx: ctx, userID: userID, playerID: playerID, lobbyID: lobbyID, deck: deck, slot: slot}
}

func (e ParticipantCreatedEvent) Name() string {
	return "lobby.participant.created"
}

func (e ParticipantCreatedEvent) ToSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.String("lobby.id", e.lobbyID.String()),
		slog.String("user.id", e.userID.String()),
		slog.String("player.id", e.playerID.String()),
		slog.Int("deck.index", e.deck),
		slog.Int("slot.index", e.slot),
	}
}

func (e ParticipantCreatedEvent) Context() context.Context {
	return e.ctx
}

func (e ParticipantCreatedEvent) LobbyID() uuid.UUID {
	return e.lobbyID
}

func (e ParticipantCreatedEvent) UserID() uuid.UUID {
	return e.userID
}

func (e ParticipantCreatedEvent) PlayerID() uuid.UUID {
	return e.playerID
}

func (e ParticipantCreatedEvent) DeckIndex() int {
	return e.deck
}

func (e ParticipantCreatedEvent) SlotIndex() int {
	return e.slot
}
