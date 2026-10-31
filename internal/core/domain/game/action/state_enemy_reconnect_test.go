package action_test

import (
	"testing"
	"time"

	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEnemyTurnGame(t *testing.T) (*game.LiveGameInstance, *game.LivePlayer) {
	t.Helper()
	g, player := newTurnGateGame(t)
	drainChanges(g)
	t.Cleanup(g.Stop)
	g.State = action.NewEnemyTurnState(g)
	return g, player
}

func TestEnemyTurnState_NoReconnect_TimeoutLiftsReadyWait(t *testing.T) {
	g, _ := newEnemyTurnGame(t)

	g.State.Update(g, time.Now().UTC().Add(time.Second*20))

	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}

func TestEnemyTurnState_Reconnect_RestartsReadyWait(t *testing.T) {
	g, player := newEnemyTurnGame(t)
	var start = time.Now().UTC()
	player.Disconnected = true

	// Away past the timeout: the disconnect hold keeps the turn.
	g.State.Update(g, start.Add(time.Second*20))
	require.IsType(t, &action.EnemyTurnState{}, g.State)

	var back = start.Add(time.Second * 20)
	require.NoError(t, action.NewPlayerReconnectAction(player.PlayerID, back).Perform(g))

	// Back but not ready: still playing the round out, so the turn waits for them.
	g.State.Update(g, back.Add(time.Second*5))
	assert.IsType(t, &action.EnemyTurnState{}, g.State)

	// A silent client still cannot stall it past a fresh timeout.
	g.State.Update(g, back.Add(time.Second*16))
	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}

func TestEnemyTurnState_Reconnect_ReadyEndsWaitEarly(t *testing.T) {
	g, player := newEnemyTurnGame(t)
	var back = time.Now().UTC().Add(time.Second * 20)
	player.Disconnected = true
	require.NoError(t, action.NewPlayerReconnectAction(player.PlayerID, back).Perform(g))

	player.Ready = true
	g.State.Update(g, back.Add(time.Second))

	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}
