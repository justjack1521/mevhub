package action_test

import (
	"testing"
	"time"

	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newAwayGame builds a two-player game at the start of a player turn: one who
// plays every round and one who can be left silent.
func newAwayGame(t *testing.T) (g *game.LiveGameInstance, active, silent *game.LivePlayer) {
	t.Helper()
	g, active = newTurnGateGame(t)
	t.Cleanup(g.Stop)
	silent = &game.LivePlayer{UserID: uuid.NewV4(), PlayerID: uuid.NewV4(), PartySlot: 1}
	for _, party := range g.Parties {
		party.Players[silent.PlayerID] = silent
	}
	g.State = action.NewPlayerTurnState(g)
	takeChanges(g)
	return g, active, silent
}

// takeChanges empties the change channel without blocking and returns what was
// in it, so a test never stalls on the channel's buffer.
func takeChanges(g *game.LiveGameInstance) []game.Change {
	var changes []game.Change
	for {
		select {
		case change := <-g.ChangeChannel:
			changes = append(changes, change)
		default:
			return changes
		}
	}
}

func countOf[T game.Change](changes []game.Change) int {
	var n = 0
	for _, change := range changes {
		if _, ok := change.(T); ok {
			n++
		}
	}
	return n
}

// playRound runs one whole round with only the active player sending anything:
// they lock, the silent player is force-locked when the turn times out, then
// the active player readies and the enemy turn times out on the silent one.
func playRound(t *testing.T, g *game.LiveGameInstance, active *game.LivePlayer) []game.Change {
	t.Helper()
	require.IsType(t, &action.PlayerTurnState{}, g.State)
	require.NoError(t, action.NewPlayerLockAction(g.InstanceID, uuid.Nil, active.PlayerID).Perform(g))
	g.State.Update(g, time.Now().UTC().Add(time.Second*31))
	require.IsType(t, &action.EnemyTurnState{}, g.State)
	var changes = takeChanges(g)
	require.NoError(t, action.NewPlayerReadyAction(g.InstanceID, uuid.Nil, active.PlayerID).Perform(g))
	g.State.Update(g, time.Now().UTC().Add(time.Second*16))
	require.IsType(t, &action.PlayerTurnState{}, g.State)
	return append(changes, takeChanges(g)...)
}

func TestAway_SilentRoundsMarkPlayerAway(t *testing.T) {
	g, active, silent := newAwayGame(t)

	var first = playRound(t, g, active)
	assert.False(t, silent.Away, "one silent round is not enough")
	assert.Equal(t, 1, silent.SilentRounds)
	assert.Zero(t, countOf[*action.PlayerAwayChange](first))

	var second = playRound(t, g, active)
	assert.True(t, silent.Away)
	assert.Equal(t, 1, countOf[*action.PlayerAwayChange](second))
	assert.False(t, active.Away)
	assert.Zero(t, active.SilentRounds)
}

func TestAway_AnyMessageResetsTheCount(t *testing.T) {
	g, active, silent := newAwayGame(t)

	playRound(t, g, active)
	require.Equal(t, 1, silent.SilentRounds)

	// A stamp is enough: the client is alive even though it locked nothing.
	require.NoError(t, action.NewPlayerStampAction(g.InstanceID, silent.PlayerID, uuid.NewV4(), time.Now().UTC()).Perform(g))
	playRound(t, g, active)

	assert.False(t, silent.Away)
	assert.Zero(t, silent.SilentRounds)
}

func TestAway_TurnsDoNotWaitOnAwayPlayer(t *testing.T) {
	g, active, silent := newAwayGame(t)
	playRound(t, g, active)
	playRound(t, g, active)
	require.True(t, silent.Away)

	// The player turn ends as soon as the active player locks, and the away
	// player is locked on the way out.
	require.NoError(t, action.NewPlayerLockAction(g.InstanceID, uuid.Nil, active.PlayerID).Perform(g))
	g.State.Update(g, time.Now().UTC().Add(time.Second))
	require.IsType(t, &action.EnemyTurnState{}, g.State)
	assert.True(t, silent.ActionsLocked)
	takeChanges(g)

	// The enemy turn ends as soon as the active player is ready.
	require.NoError(t, action.NewPlayerReadyAction(g.InstanceID, uuid.Nil, active.PlayerID).Perform(g))
	g.State.Update(g, time.Now().UTC().Add(time.Second))
	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}

func TestAway_ActivityReturnsPlayerAndTurnWaitsAgain(t *testing.T) {
	g, active, silent := newAwayGame(t)
	playRound(t, g, active)
	playRound(t, g, active)
	require.True(t, silent.Away)

	require.NoError(t, newEnqueue(g, silent).Perform(g))
	assert.False(t, silent.Away)
	assert.Equal(t, 1, countOf[*action.PlayerReturnChange](takeChanges(g)))

	// Back, so the turn waits for their lock again.
	require.NoError(t, action.NewPlayerLockAction(g.InstanceID, uuid.Nil, active.PlayerID).Perform(g))
	g.State.Update(g, time.Now().UTC().Add(time.Second))
	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}

func TestAway_EveryoneAwayFallsBackToTimers(t *testing.T) {
	g, active, silent := newAwayGame(t)
	active.Away = true
	silent.Away = true

	// Nobody present must not read as everyone locked: the turn runs its timer.
	g.State.Update(g, time.Now().UTC().Add(time.Second))
	require.IsType(t, &action.PlayerTurnState{}, g.State)

	g.State.Update(g, time.Now().UTC().Add(time.Second*31))
	require.IsType(t, &action.EnemyTurnState{}, g.State)
	takeChanges(g)

	g.State.Update(g, time.Now().UTC().Add(time.Second))
	require.IsType(t, &action.EnemyTurnState{}, g.State)

	g.State.Update(g, time.Now().UTC().Add(time.Second*16))
	assert.IsType(t, &action.PlayerTurnState{}, g.State)
}

func TestAway_ReconnectCountsAsActivity(t *testing.T) {
	g, active, silent := newAwayGame(t)
	silent.Disconnected = true
	silent.DisconnectTime = time.Now().UTC()

	// Dropped and back within the round, then nothing else: the reconnect alone
	// shows a live client, so the round is not held against them.
	require.NoError(t, action.NewPlayerReconnectAction(silent.PlayerID, time.Now().UTC()).Perform(g))
	playRound(t, g, active)

	assert.Zero(t, silent.SilentRounds)
	assert.False(t, silent.Away)
}

func TestAway_ReconnectClearsAway(t *testing.T) {
	g, _, silent := newAwayGame(t)
	silent.Away = true
	silent.SilentRounds = 3
	silent.Disconnected = true
	silent.DisconnectTime = time.Now().UTC()

	require.NoError(t, action.NewPlayerReconnectAction(silent.PlayerID, time.Now().UTC()).Perform(g))

	assert.False(t, silent.Away)
	assert.Zero(t, silent.SilentRounds)
	// The reconnect notification clears the badge; a return on top would be noise.
	assert.Zero(t, countOf[*action.PlayerReturnChange](takeChanges(g)))
}

func TestGameSync_AwayShowsAsDisconnected(t *testing.T) {
	g, _, silent := newAwayGame(t)
	silent.Away = true

	var sync = action.NewGameStateSyncChange(g)
	for _, party := range sync.Parties {
		for _, player := range party.Players {
			assert.Equal(t, player.PlayerID == silent.PlayerID, player.Disconnected)
		}
	}
}
