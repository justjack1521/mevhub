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

// performStamp runs a stamp and returns the change it broadcast, or nil when it
// was dropped.
func performStamp(t *testing.T, g *game.LiveGameInstance, player uuid.UUID, stamp uuid.UUID, at time.Time) *action.PlayerStampChange {
	t.Helper()
	require.NoError(t, action.NewPlayerStampAction(g.InstanceID, player, stamp, at).Perform(g))
	select {
	case change := <-g.ChangeChannel:
		actual, ok := change.(*action.PlayerStampChange)
		require.True(t, ok, "unexpected change %T", change)
		return actual
	default:
		return nil
	}
}

func TestPlayerStamp_Broadcasts(t *testing.T) {
	// The claim fixture's player is dead, which is the point: the fallen may still stamp.
	g, player := newClaimGame(t)
	stamp := uuid.NewV4()

	change := performStamp(t, g, player.PlayerID, stamp, time.Now().UTC())

	require.NotNil(t, change)
	assert.Equal(t, stamp, change.StampID)
	assert.Equal(t, player.PartySlot, change.PartySlot)
}

func TestPlayerStamp_InsideCooldown_Dropped(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()

	require.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now))
	assert.Nil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now.Add(game.StampCooldown-time.Millisecond)))
}

func TestPlayerStamp_AfterCooldown_Broadcasts(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()

	require.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now))
	assert.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now.Add(game.StampCooldown)))
}

func TestPlayerStamp_DroppedStamp_DoesNotRestartCooldown(t *testing.T) {
	g, player := newClaimGame(t)
	now := time.Now().UTC()

	require.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now))
	require.Nil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now.Add(time.Second)))
	assert.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now.Add(game.StampCooldown)))
}

func TestPlayerStamp_CooldownIsPerPlayer(t *testing.T) {
	g, player := newClaimGame(t)
	var ally *game.LivePlayer
	for _, party := range g.Parties {
		for _, p := range party.Players {
			if uuid.Equal(p.PlayerID, player.PlayerID) == false {
				ally = p
			}
		}
	}
	now := time.Now().UTC()

	require.NotNil(t, performStamp(t, g, player.PlayerID, uuid.NewV4(), now))
	assert.NotNil(t, performStamp(t, g, ally.PlayerID, uuid.NewV4(), now))
}

func TestPlayerStamp_PlayerNotInGame_Errors(t *testing.T) {
	g, _ := newClaimGame(t)
	assert.Error(t, action.NewPlayerStampAction(g.InstanceID, uuid.NewV4(), uuid.NewV4(), time.Now().UTC()).Perform(g))
}
