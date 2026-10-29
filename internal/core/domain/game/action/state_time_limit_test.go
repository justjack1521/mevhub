package action_test

import (
	"testing"
	"time"

	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// drainChanges empties the change channel so a state transition, which sends
// on it, cannot block the test.
func drainChanges(g *game.LiveGameInstance) {
	go func() {
		for {
			select {
			case <-g.ChangeChannel:
			case <-g.Done():
				return
			}
		}
	}()
}

func newTimeLimitGame(t *testing.T) *game.LiveGameInstance {
	t.Helper()
	g, _ := newTurnGateGame(t)
	drainChanges(g)
	t.Cleanup(g.Stop)
	return g
}

func TestTimeExpired_NotStarted_NeverExpires(t *testing.T) {
	g := newTimeLimitGame(t)

	assert.False(t, g.TimeExpired(time.Now().Add(time.Hour)))
}

func TestTimeExpired_NoLimit_NeverExpires(t *testing.T) {
	g := newTimeLimitGame(t)
	g.GameDuration = 0
	g.StartedAt = time.Now()

	assert.False(t, g.TimeExpired(g.StartedAt.Add(time.Hour)))
}

func TestTimeExpired_AtLimit_Expires(t *testing.T) {
	g := newTimeLimitGame(t)
	g.StartedAt = time.Now()

	assert.False(t, g.TimeExpired(g.StartedAt.Add(g.GameDuration-time.Second)))
	assert.True(t, g.TimeExpired(g.StartedAt.Add(g.GameDuration)))
}

func TestPendingState_AllReady_StartsClock(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewPendingState(g)
	for _, party := range g.Parties {
		for _, player := range party.Players {
			player.Ready = true
		}
	}
	now := time.Now().UTC()

	g.State.Update(g, now)

	require.IsType(t, &action.PlayerTurnState{}, g.State)
	assert.Equal(t, now, g.StartedAt)
}

func TestPlayerTurnState_TimeExpired_EndsGame(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewPlayerTurnState(g)
	g.StartedAt = time.Now().UTC().Add(-g.GameDuration)

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.EndGameState{}, g.State)
	assert.True(t, g.HasEnded())
}

func TestPlayerTurnState_TimeRemaining_KeepsTurn(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewPlayerTurnState(g)
	g.StartedAt = time.Now().UTC()

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.PlayerTurnState{}, g.State)
	assert.False(t, g.HasEnded())
}

func TestEnemyTurnState_TimeExpired_WaitsForBoundary(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewEnemyTurnState(g)
	g.StartedAt = time.Now().UTC().Add(-g.GameDuration)

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.EnemyTurnState{}, g.State)
	assert.False(t, g.HasEnded())
}

func TestEnemyTurnState_TimeExpiredAtBoundary_EndsGame(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewEnemyTurnState(g)
	g.StartedAt = time.Now().UTC().Add(-g.GameDuration)
	for _, party := range g.Parties {
		for _, player := range party.Players {
			player.Ready = true
		}
	}

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.EndGameState{}, g.State)
	assert.True(t, g.HasEnded())
}

func TestEnemyTurnState_TimeRemainingAtBoundary_OpensPlayerTurn(t *testing.T) {
	g := newTimeLimitGame(t)
	g.State = action.NewEnemyTurnState(g)
	g.StartedAt = time.Now().UTC()
	for _, party := range g.Parties {
		for _, player := range party.Players {
			player.Ready = true
		}
	}

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.PlayerTurnState{}, g.State)
	assert.False(t, g.HasEnded())
}
