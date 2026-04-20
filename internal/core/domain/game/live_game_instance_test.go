package game_test

import (
	"testing"
	"time"

	"mevhub/internal/core/domain/game"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- helpers ---

func newLivePlayer(maxActions int) *game.LivePlayer {
	return &game.LivePlayer{
		UserID:         uuid.NewV4(),
		PlayerID:       uuid.NewV4(),
		MaxActionCount: maxActions,
	}
}

func newLiveParty(maxPlayers int) *game.LiveParty {
	return &game.LiveParty{
		PartyID:            uuid.NewV4(),
		Players:            make(map[uuid.UUID]*game.LivePlayer),
		MaxPlayerCount:     maxPlayers,
		PlayerTurnDuration: time.Second * 30,
	}
}

func newLiveGame(t *testing.T) *game.LiveGameInstance {
	t.Helper()
	inst := &game.Instance{
		SysID: uuid.NewV4(),
		Options: &game.InstanceOptions{
			MaxRunTime:         time.Minute * 10,
			MaxPartyCount:      2,
			MaxPlayerCount:     4,
			PlayerTurnDuration: time.Second * 30,
		},
	}
	return game.NewLiveGameInstance(inst)
}

// --- LivePlayer: EnqueueAction ---

func TestEnqueueAction_BelowMax_Succeeds(t *testing.T) {
	p := newLivePlayer(3)
	action := &game.PlayerAction{ActionType: game.PlayerActionTypeNormalAttack}
	err := p.EnqueueAction(action)
	assert.NoError(t, err)
	assert.Len(t, p.Actions, 1)
}

func TestEnqueueAction_AtMaxCapacity_ReturnsError(t *testing.T) {
	p := newLivePlayer(2)
	_ = p.EnqueueAction(&game.PlayerAction{})
	_ = p.EnqueueAction(&game.PlayerAction{})
	err := p.EnqueueAction(&game.PlayerAction{})
	assert.ErrorIs(t, err, game.ErrPlayerActionsFull)
}

func TestEnqueueAction_ZeroMaxCapacity_UnlimitedQueue(t *testing.T) {
	p := newLivePlayer(0) // 0 = unlimited
	for i := 0; i < 10; i++ {
		err := p.EnqueueAction(&game.PlayerAction{})
		assert.NoError(t, err)
	}
	assert.Len(t, p.Actions, 10)
}

func TestEnqueueAction_ActionsLocked_ReturnsError(t *testing.T) {
	p := newLivePlayer(5)
	p.ActionsLocked = true
	err := p.EnqueueAction(&game.PlayerAction{})
	assert.ErrorIs(t, err, game.ErrPlayerActionsLocked)
}

// --- LivePlayer: DequeueAction ---

func TestDequeueAction_WithActions_Succeeds(t *testing.T) {
	p := newLivePlayer(3)
	_ = p.EnqueueAction(&game.PlayerAction{ActionType: game.PlayerActionTypeNormalAttack})
	_ = p.EnqueueAction(&game.PlayerAction{ActionType: game.PlayerActionTypeAbilityCast})
	err := p.DequeueAction()
	assert.NoError(t, err)
	assert.Len(t, p.Actions, 1)
}

func TestDequeueAction_EmptyQueue_ReturnsError(t *testing.T) {
	p := newLivePlayer(3)
	err := p.DequeueAction()
	assert.ErrorIs(t, err, game.ErrPlayerActionsEmpty)
}

func TestDequeueAction_ActionsLocked_ReturnsError(t *testing.T) {
	p := newLivePlayer(3)
	_ = p.EnqueueAction(&game.PlayerAction{})
	p.ActionsLocked = true
	err := p.DequeueAction()
	assert.ErrorIs(t, err, game.ErrPlayerActionsLocked)
}

// --- LiveParty ---

func TestLiveParty_PlayerExists_ReturnsFalseForUnknown(t *testing.T) {
	party := newLiveParty(4)
	assert.False(t, party.PlayerExists(uuid.NewV4()))
}

func TestLiveParty_GetPlayer_UnknownID_ReturnsError(t *testing.T) {
	party := newLiveParty(4)
	_, err := party.GetPlayer(uuid.NewV4())
	assert.Error(t, err)
}

func TestLiveParty_GetReadyPlayerCount_CountsCorrectly(t *testing.T) {
	party := newLiveParty(4)
	p1 := newLivePlayer(3)
	p1.Ready = true
	p2 := newLivePlayer(3)
	p2.Ready = false
	party.Players[p1.PlayerID] = p1
	party.Players[p2.PlayerID] = p2

	assert.Equal(t, 1, party.GetReadyPlayerCount())
}

func TestLiveParty_GetActionLockedPlayerCount_CountsCorrectly(t *testing.T) {
	party := newLiveParty(4)
	p1 := newLivePlayer(3)
	p1.ActionsLocked = true
	p2 := newLivePlayer(3)
	p2.ActionsLocked = true
	p3 := newLivePlayer(3)
	p3.ActionsLocked = false
	party.Players[p1.PlayerID] = p1
	party.Players[p2.PlayerID] = p2
	party.Players[p3.PlayerID] = p3

	assert.Equal(t, 2, party.GetActionLockedPlayerCount())
}

func TestLiveParty_RemovePlayer_Succeeds(t *testing.T) {
	party := newLiveParty(4)
	p := newLivePlayer(3)
	party.Players[p.PlayerID] = p

	err := party.RemovePlayer(p.PlayerID)
	assert.NoError(t, err)
	assert.False(t, party.PlayerExists(p.PlayerID))
}

func TestLiveParty_RemovePlayer_UnknownID_ReturnsError(t *testing.T) {
	party := newLiveParty(4)
	err := party.RemovePlayer(uuid.NewV4())
	assert.Error(t, err)
}

// --- LiveGameInstance ---

func TestLiveGameInstance_GetPlayer_AcrossParties(t *testing.T) {
	g := newLiveGame(t)

	party := newLiveParty(4)
	player := newLivePlayer(3)
	party.Players[player.PlayerID] = player
	g.Parties[party.PartyID] = party

	found, err := g.GetPlayer(player.PlayerID)
	require.NoError(t, err)
	assert.Equal(t, player.PlayerID, found.PlayerID)
}

func TestLiveGameInstance_GetPlayer_NotFound_ReturnsError(t *testing.T) {
	g := newLiveGame(t)
	_, err := g.GetPlayer(uuid.NewV4())
	assert.Error(t, err)
}

func TestLiveGameInstance_RemovePlayer_Succeeds(t *testing.T) {
	g := newLiveGame(t)
	party := newLiveParty(4)
	player := newLivePlayer(3)
	party.Players[player.PlayerID] = player
	g.Parties[party.PartyID] = party

	err := g.RemovePlayer(player.PlayerID)
	assert.NoError(t, err)
	assert.False(t, g.PlayerExists(player.PlayerID))
}

func TestLiveGameInstance_RemovePlayer_NotFound_ReturnsError(t *testing.T) {
	g := newLiveGame(t)
	err := g.RemovePlayer(uuid.NewV4())
	assert.Error(t, err)
}

func TestLiveGameInstance_GetReadyPlayerCount_SumsAcrossParties(t *testing.T) {
	g := newLiveGame(t)

	for i := 0; i < 2; i++ {
		party := newLiveParty(4)
		for j := 0; j < 2; j++ {
			p := newLivePlayer(3)
			p.Ready = j == 0 // first player in each party is ready
			party.Players[p.PlayerID] = p
		}
		g.Parties[party.PartyID] = party
	}

	assert.Equal(t, 2, g.GetReadyPlayerCount())
}
