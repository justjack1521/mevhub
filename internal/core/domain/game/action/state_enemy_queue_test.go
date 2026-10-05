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

func newQueueGame(t *testing.T, count int) (*game.LiveGameInstance, *game.LiveParty, []*game.LivePlayer) {
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
	players := make([]*game.LivePlayer, count)
	for i := 0; i < count; i++ {
		players[i] = &game.LivePlayer{UserID: uuid.NewV4(), PlayerID: uuid.NewV4(), PartySlot: i}
		party.Players[players[i].PlayerID] = players[i]
	}
	g.Parties[party.PartyID] = party
	return g, party, players
}

func lock(t *testing.T, g *game.LiveGameInstance, party *game.LiveParty, player *game.LivePlayer) {
	t.Helper()
	require.NoError(t, action.NewPlayerLockAction(g.InstanceID, party.PartyID, player.PlayerID).Perform(g))
}

func queuedIDs(state *action.EnemyTurnState, party *game.LiveParty) []uuid.UUID {
	var ids []uuid.UUID
	for _, q := range state.QueuedActions[party.PartyIndex] {
		ids = append(ids, q.PlayerID)
	}
	return ids
}

func TestLockIndex_PlayerRemovedAfterLocking_NoCollision(t *testing.T) {
	g, party, p := newQueueGame(t, 3)
	g.State = action.NewPlayerTurnState(g)

	lock(t, g, party, p[0])
	lock(t, g, party, p[1])
	require.NoError(t, party.RemovePlayer(p[0].PlayerID))
	lock(t, g, party, p[2])

	assert.NotEqual(t, p[1].ActionLockIndex, p[2].ActionLockIndex)

	state := action.NewEnemyTurnState(g)
	assert.Equal(t, []uuid.UUID{p[1].PlayerID, p[2].PlayerID}, queuedIDs(state, party))
}

func TestEnemyTurnQueue_DeadPlayer_Packed(t *testing.T) {
	g, party, p := newQueueGame(t, 3)
	g.State = action.NewPlayerTurnState(g)

	lock(t, g, party, p[2])
	lock(t, g, party, p[0])
	p[1].Dead = true

	state := action.NewEnemyTurnState(g)
	assert.Equal(t, []uuid.UUID{p[2].PlayerID, p[0].PlayerID}, queuedIDs(state, party))
}
