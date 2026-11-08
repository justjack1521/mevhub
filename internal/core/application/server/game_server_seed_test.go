package server

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSeedTestServer(t *testing.T) *GameServer {
	t.Helper()
	var svr = NewGameServerFactory(nil).Create(&game.Instance{
		SysID: uuid.NewV4(),
		Options: &game.InstanceOptions{
			MaxRunTime:         time.Minute * 10,
			MaxPartyCount:      1,
			MaxPlayerCount:     4,
			PlayerTurnDuration: time.Second * 30,
		},
	})
	t.Cleanup(svr.Stop)
	return svr
}

func TestGameServer_Seed_PopulatesBeforeStart(t *testing.T) {
	var svr = newSeedTestServer(t)
	var party, user, player = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

	require.NoError(t, svr.Seed(action.NewPartyAddAction(party, 0)))
	require.NoError(t, svr.Seed(action.NewPlayerAddAction(user, player, party, 1)))

	assert.Equal(t, 1, svr.game.GetPlayerCount(), "the player must exist before the loop can receive any action for them")
	assert.Len(t, svr.game.ChangeChannel, 2, "seeded changes wait for Start to be handled")
	assert.Empty(t, svr.clients, "the client registry is filled from the changes, once Start runs")
}

func TestGameServer_Seed_RefusesWhenChangeChannelFull(t *testing.T) {
	var svr = newSeedTestServer(t)
	for len(svr.game.ChangeChannel) < cap(svr.game.ChangeChannel) {
		svr.game.ChangeChannel <- action.NewPartyAddChange(uuid.NewV4(), 0)
	}

	var done = make(chan error, 1)
	go func() { done <- svr.Seed(action.NewPartyAddAction(uuid.NewV4(), 0)) }()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, ErrGameServerSeedOverflow)
	case <-time.After(time.Second):
		t.Fatal("Seed blocked on a full change channel")
	}
}

func TestGameServerHost_Action_FullGameQueueDoesNotBlock(t *testing.T) {
	var host = NewGameServerHost(slog.New(slog.NewTextHandler(io.Discard, nil)), NewGameServerFactory(nil), nil)
	var svr = newSeedTestServer(t)
	host.games[svr.InstanceID] = svr
	// Never started: nothing drains the game's queue, as with a stalled loop.
	for len(svr.game.ActionChannel) < cap(svr.game.ActionChannel) {
		svr.game.ActionChannel <- nil
	}

	var done = make(chan struct{})
	go func() {
		host.action(&GameActionRequest{GameID: svr.InstanceID, Action: action.NewPartyAddAction(uuid.NewV4(), 0)})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("routing to a game with a full queue blocked the host")
	}
}
