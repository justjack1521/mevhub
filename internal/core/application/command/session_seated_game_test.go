package command

import (
	"context"
	"errors"
	"testing"

	"mevhub/internal/core/domain/session"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type seatedGameContext struct {
	context.Context
	user, player uuid.UUID
}

func (c seatedGameContext) UserID() uuid.UUID   { return c.user }
func (c seatedGameContext) PlayerID() uuid.UUID { return c.player }

type seatedGameSessions struct {
	updates []session.Instance
	err     error
}

func (r *seatedGameSessions) Create(context.Context, *session.Instance) error { return nil }
func (r *seatedGameSessions) Update(_ context.Context, instance *session.Instance) error {
	if r.err != nil {
		return r.err
	}
	r.updates = append(r.updates, *instance)
	return nil
}
func (r *seatedGameSessions) UpdateConnectionState(context.Context, *session.Instance) error {
	return nil
}
func (r *seatedGameSessions) Delete(context.Context, *session.Instance) error { return nil }

type seatedGameRecorder struct {
	events []session.GameAbandonedEvent
}

func (r *seatedGameRecorder) Notify(event mevent.Event) {
	r.events = append(r.events, event.(session.GameAbandonedEvent))
}

func newSeatedGameFixture() (seatedGameContext, *mevent.Publisher, *seatedGameRecorder, *seatedGameSessions) {
	var ctx = seatedGameContext{Context: context.Background(), user: uuid.NewV4(), player: uuid.NewV4()}
	var publisher = mevent.NewPublisher()
	var recorder = &seatedGameRecorder{}
	publisher.Subscribe(recorder, session.GameAbandonedEvent{})
	return ctx, publisher, recorder, &seatedGameSessions{}
}

func TestLeaveSeatedGame_NoGameDoesNothing(t *testing.T) {
	ctx, publisher, recorder, sessions := newSeatedGameFixture()
	var current = &session.Instance{UserID: ctx.user, PlayerID: ctx.player}

	if err := leaveSeatedGame(ctx, publisher, sessions, current); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessions.updates) != 0 {
		t.Errorf("session updated %d times, want 0", len(sessions.updates))
	}
	if len(recorder.events) != 0 {
		t.Errorf("published %d events, want 0", len(recorder.events))
	}
}

func TestLeaveSeatedGame_ClearsGameAndPublishesRemoval(t *testing.T) {
	ctx, publisher, recorder, sessions := newSeatedGameFixture()
	var seated = uuid.NewV4()
	var current = &session.Instance{UserID: ctx.user, PlayerID: ctx.player, GameID: seated}

	if err := leaveSeatedGame(ctx, publisher, sessions, current); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !uuid.Equal(current.GameID, uuid.Nil) {
		t.Errorf("session still names game %s", current.GameID)
	}
	if len(sessions.updates) != 1 || !uuid.Equal(sessions.updates[0].GameID, uuid.Nil) {
		t.Fatalf("persisted %+v, want one update with no game", sessions.updates)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("published %d events, want 1", len(recorder.events))
	}
	var event = recorder.events[0]
	if !uuid.Equal(event.GameID(), seated) || !uuid.Equal(event.UserID(), ctx.user) || !uuid.Equal(event.PlayerID(), ctx.player) {
		t.Errorf("event names game %s user %s player %s, want %s %s %s", event.GameID(), event.UserID(), event.PlayerID(), seated, ctx.user, ctx.player)
	}
}

func TestLeaveSeatedGame_FailedUpdateKeepsGameAndPublishesNothing(t *testing.T) {
	ctx, publisher, recorder, sessions := newSeatedGameFixture()
	sessions.err = errors.New("redis down")
	var seated = uuid.NewV4()
	var current = &session.Instance{UserID: ctx.user, PlayerID: ctx.player, GameID: seated}

	if err := leaveSeatedGame(ctx, publisher, sessions, current); err == nil {
		t.Fatal("expected the update error")
	}
	if !uuid.Equal(current.GameID, seated) {
		t.Errorf("session game is %s, want it left as %s", current.GameID, seated)
	}
	if len(recorder.events) != 0 {
		t.Errorf("published %d events, want 0", len(recorder.events))
	}
}
