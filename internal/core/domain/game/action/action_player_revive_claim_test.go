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

func newClaimGame(t *testing.T) (*game.LiveGameInstance, *game.LivePlayer) {
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
		UserID:    uuid.NewV4(),
		PlayerID:  uuid.NewV4(),
		Dead:      true,
		DeathTime: time.Now().UTC(),
	}
	// A standing teammate: reviving is only possible while someone is alive.
	// Without one the party is wiped, which newWipedGame covers.
	ally := &game.LivePlayer{
		UserID:    uuid.NewV4(),
		PlayerID:  uuid.NewV4(),
		PartySlot: 1,
	}
	party.Players[player.PlayerID] = player
	party.Players[ally.PlayerID] = ally
	g.Parties[party.PartyID] = party
	return g, player
}

// newWipedGame is newClaimGame with every player down.
func newWipedGame(t *testing.T) (*game.LiveGameInstance, *game.LivePlayer) {
	t.Helper()
	g, player := newClaimGame(t)
	for _, party := range g.Parties {
		for _, p := range party.Players {
			p.Dead = true
			p.DeathTime = time.Now().UTC()
		}
	}
	return g, player
}

func performClaim(t *testing.T, g *game.LiveGameInstance, target uuid.UUID, source uuid.UUID, at time.Time) action.PlayerReviveClaimResult {
	t.Helper()
	response := make(chan action.PlayerReviveClaimResult, 1)
	claim := action.NewPlayerReviveClaimAction(g.InstanceID, target, source, at, response)
	require.NoError(t, claim.Perform(g))
	select {
	case result := <-response:
		return result
	default:
		t.Fatal("claim produced no reply")
		return action.PlayerReviveClaimResult{}
	}
}

func TestPlayerReviveClaim_DeadUnclaimed_Granted(t *testing.T) {
	g, player := newClaimGame(t)
	source := uuid.NewV4()
	now := time.Now().UTC()

	result := performClaim(t, g, player.PlayerID, source, now)

	assert.True(t, result.Granted)
	assert.Equal(t, action.ReviveClaimDenyReasonNone, result.Reason)
	assert.Equal(t, game.ReviveClaimDuration, result.Remaining)
	assert.Equal(t, source, player.ReviveClaimSource)
	assert.Equal(t, now.Add(game.ReviveClaimDuration), player.ReviveClaimExpiry)

	select {
	case change := <-g.ChangeChannel:
		claim, ok := change.(*action.PlayerReviveClaimChange)
		require.True(t, ok, "expected PlayerReviveClaimChange, got %T", change)
		assert.Equal(t, player.PlayerID, claim.PlayerID)
		assert.Equal(t, source, claim.SourceID)
		assert.Equal(t, game.ReviveClaimDuration, claim.Remaining)
	default:
		t.Fatal("grant broadcast no change")
	}
}

func TestPlayerReviveClaim_TargetAlive_Denied(t *testing.T) {
	g, player := newClaimGame(t)
	player.Dead = false

	result := performClaim(t, g, player.PlayerID, uuid.NewV4(), time.Now().UTC())

	assert.False(t, result.Granted)
	assert.Equal(t, action.ReviveClaimDenyReasonTargetAlive, result.Reason)
	assert.Equal(t, uuid.Nil, player.ReviveClaimSource)
}

func TestPlayerReviveClaim_TargetNotFound_Denied(t *testing.T) {
	g, _ := newClaimGame(t)

	result := performClaim(t, g, uuid.NewV4(), uuid.NewV4(), time.Now().UTC())

	assert.False(t, result.Granted)
	assert.Equal(t, action.ReviveClaimDenyReasonTargetNotFound, result.Reason)
}

func TestPlayerReviveClaim_HeldByOther_DeniedWithRemaining(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()
	first := uuid.NewV4()
	performClaim(t, g, player.PlayerID, first, now)

	later := now.Add(time.Second * 4)
	result := performClaim(t, g, player.PlayerID, uuid.NewV4(), later)

	assert.False(t, result.Granted)
	assert.Equal(t, action.ReviveClaimDenyReasonAlreadyClaimed, result.Reason)
	assert.Equal(t, game.ReviveClaimDuration-time.Second*4, result.Remaining)
	assert.Equal(t, first, player.ReviveClaimSource)
}

func TestPlayerReviveClaim_SameClaimantRetry_RefreshesWindow(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()
	source := uuid.NewV4()
	performClaim(t, g, player.PlayerID, source, now)

	later := now.Add(time.Second * 4)
	result := performClaim(t, g, player.PlayerID, source, later)

	assert.True(t, result.Granted)
	assert.Equal(t, later.Add(game.ReviveClaimDuration), player.ReviveClaimExpiry)
}

func TestPlayerReviveClaim_ExpiredClaim_Regrantable(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()
	performClaim(t, g, player.PlayerID, uuid.NewV4(), now)

	second := uuid.NewV4()
	afterDecay := now.Add(game.ReviveClaimDuration + time.Millisecond)
	result := performClaim(t, g, player.PlayerID, second, afterDecay)

	assert.True(t, result.Granted)
	assert.Equal(t, second, player.ReviveClaimSource)
}

func TestPlayerRevive_ClearsStandingClaim(t *testing.T) {
	g, player := newClaimGame(t)
	source := uuid.NewV4()
	performClaim(t, g, player.PlayerID, source, time.Now().UTC())

	revive := action.NewPlayerReviveAction(g.InstanceID, player.PlayerID, source)
	require.NoError(t, revive.Perform(g))

	assert.False(t, player.Dead)
	assert.Equal(t, uuid.Nil, player.ReviveClaimSource)
	assert.True(t, player.ReviveClaimExpiry.IsZero())
}

func TestPlayerReviveClaim_PartyWiped_Denied(t *testing.T) {
	g, player := newWipedGame(t)

	result := performClaim(t, g, player.PlayerID, player.PlayerID, time.Now().UTC())

	assert.False(t, result.Granted)
	assert.Equal(t, action.ReviveClaimDenyReasonUnavailable, result.Reason)
	assert.Equal(t, uuid.Nil, player.ReviveClaimSource, "a wiped party must not take a claim")
}

func TestPlayerRevive_PartyWiped_NoOp(t *testing.T) {
	g, player := newWipedGame(t)

	revive := action.NewPlayerReviveAction(g.InstanceID, player.PlayerID, player.PlayerID)
	require.NoError(t, revive.Perform(g))

	assert.True(t, player.Dead, "a wipe is final")
}

func TestPlayerTurn_PartyWiped_EndsGame(t *testing.T) {
	g, _ := newWipedGame(t)
	g.State = action.NewPlayerTurnState(g)

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.EndGameState{}, g.State)
}

func TestEnemyTurn_PartyWiped_EndsGame(t *testing.T) {
	g, _ := newWipedGame(t)
	g.State = action.NewEnemyTurnState(g)

	g.State.Update(g, time.Now().UTC())

	assert.IsType(t, &action.EndGameState{}, g.State)
}
