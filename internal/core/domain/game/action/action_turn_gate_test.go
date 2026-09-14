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

// newTurnGateGame builds a one-party game holding a single living, unlocked
// player: the shape in which the per-player checks all pass, so any rejection
// can only come from the state gate.
func newTurnGateGame(t *testing.T) (*game.LiveGameInstance, *game.LivePlayer) {
	t.Helper()
	g := game.NewLiveGameInstance(&game.Instance{
		SysID: uuid.NewV4(),
		Options: &game.InstanceOptions{
			MaxRunTime:         time.Minute * 10,
			MaxPartyCount:      1,
			MaxPlayerCount:     4,
			PlayerTurnDuration: time.Second * 30,
		},
	})
	party := &game.LiveParty{
		PartyID:        uuid.NewV4(),
		Players:        make(map[uuid.UUID]*game.LivePlayer),
		MaxPlayerCount: 4,
	}
	player := &game.LivePlayer{
		UserID:   uuid.NewV4(),
		PlayerID: uuid.NewV4(),
	}
	party.Players[player.PlayerID] = player
	g.Parties[party.PartyID] = party
	return g, player
}

func newEnqueue(g *game.LiveGameInstance, player *game.LivePlayer) *action.PlayerEnqueueAction {
	return action.NewPlayerEnqueueAction(g.InstanceID, uuid.Nil, player.PlayerID, 0, game.PlayerActionTypeNormalAttack, 0, uuid.Nil)
}

func assertNoChange(t *testing.T, g *game.LiveGameInstance) {
	t.Helper()
	select {
	case change := <-g.ChangeChannel:
		t.Fatalf("rejected action must broadcast nothing, got %T", change)
	default:
	}
}

func TestTurnGate_Pending_EnqueueRejected(t *testing.T) {
	g, player := newTurnGateGame(t)
	g.State = action.NewPendingState(g)

	err := newEnqueue(g, player).Perform(g)

	require.ErrorIs(t, err, action.ErrNotPlayerTurn)
	assert.Empty(t, player.Actions)
	assertNoChange(t, g)
}

func TestTurnGate_Pending_LockRejected(t *testing.T) {
	g, player := newTurnGateGame(t)
	g.State = action.NewPendingState(g)

	err := action.NewPlayerLockAction(g.InstanceID, uuid.Nil, player.PlayerID).Perform(g)

	require.ErrorIs(t, err, action.ErrNotPlayerTurn)
	assert.False(t, player.ActionsLocked)
	assertNoChange(t, g)
}

func TestTurnGate_Pending_DequeueRejected(t *testing.T) {
	g, player := newTurnGateGame(t)
	g.State = action.NewPendingState(g)
	player.Actions = []*game.PlayerAction{{ActionType: game.PlayerActionTypeNormalAttack}}

	err := action.NewPlayerDequeueAction(g.InstanceID, uuid.Nil, player.PlayerID).Perform(g)

	require.ErrorIs(t, err, action.ErrNotPlayerTurn)
	assert.Len(t, player.Actions, 1)
	assertNoChange(t, g)
}

// A player who is unlocked during the enemy turn (revived mid-turn, or added
// mid-game) passes the per-player checks; the state gate must still hold.
func TestTurnGate_EnemyTurn_UnlockedPlayerRejected(t *testing.T) {
	g, player := newTurnGateGame(t)
	g.State = action.NewEnemyTurnState(g)
	require.False(t, player.ActionsLocked)

	assert.ErrorIs(t, newEnqueue(g, player).Perform(g), action.ErrNotPlayerTurn)
	assert.ErrorIs(t, action.NewPlayerLockAction(g.InstanceID, uuid.Nil, player.PlayerID).Perform(g), action.ErrNotPlayerTurn)
	assert.Empty(t, player.Actions)
	assert.False(t, player.ActionsLocked)
	assertNoChange(t, g)
}

func TestTurnGate_PlayerTurn_EnqueueLockDequeueAccepted(t *testing.T) {
	g, player := newTurnGateGame(t)
	g.State = action.NewPlayerTurnState(g)

	require.NoError(t, newEnqueue(g, player).Perform(g))
	assert.Len(t, player.Actions, 1)
	_, ok := (<-g.ChangeChannel).(*action.PlayerEnqueueActionChange)
	assert.True(t, ok)

	require.NoError(t, action.NewPlayerDequeueAction(g.InstanceID, uuid.Nil, player.PlayerID).Perform(g))
	assert.Empty(t, player.Actions)
	_, ok = (<-g.ChangeChannel).(*action.PlayerDequeueActionChange)
	assert.True(t, ok)

	require.NoError(t, action.NewPlayerLockAction(g.InstanceID, uuid.Nil, player.PlayerID).Perform(g))
	assert.True(t, player.ActionsLocked)
	_, ok = (<-g.ChangeChannel).(*action.PlayerLockActionChange)
	assert.True(t, ok)
}
