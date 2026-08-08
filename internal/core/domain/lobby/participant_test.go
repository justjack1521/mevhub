package lobby_test

import (
	"testing"

	"mevhub/internal/core/domain/lobby"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// helper: a participant with a known player ID at slot 1 (non-host)
func newParticipant(playerID uuid.UUID, slot int) *lobby.Participant {
	roleID := uuid.Nil
	return &lobby.Participant{
		UserID:          uuid.NewV4(),
		PlayerID:        playerID,
		LobbyID:         uuid.NewV4(),
		Role:            roleID,
		RoleRestriction: roleID, // no role restriction by default
		PlayerSlot:      slot,
		DeckIndex:       0,
		UseStamina:      slot == 0,
	}
}

// --- SetReady ---

func TestSetReady_CorrectPlayer_Succeeds(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	err := p.SetReady(playerID, true)
	assert.NoError(t, err)
	assert.True(t, p.Ready)
}

func TestSetReady_WrongPlayer_ReturnsError(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	err := p.SetReady(uuid.NewV4(), true)
	assert.Error(t, err)
}

// --- SetUseStamina ---

func TestSetUseStamina_NonHostCanDisable_Succeeds(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1) // slot 1 = not host
	err := p.SetUseStamina(playerID, false)
	assert.NoError(t, err)
	assert.False(t, p.UseStamina)
}

func TestSetUseStamina_HostCannotDisable_ReturnsError(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 0) // slot 0 = host
	err := p.SetUseStamina(playerID, false)
	assert.Error(t, err)
}

func TestSetUseStamina_WrongPlayer_ReturnsError(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	err := p.SetUseStamina(uuid.NewV4(), false)
	assert.Error(t, err)
}

// --- SetDeckIndex ---

func TestSetDeckIndex_NonNegative_Succeeds(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	err := p.SetDeckIndex(playerID, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, p.DeckIndex)
}

func TestSetDeckIndex_Negative_ReturnsError(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	err := p.SetDeckIndex(playerID, -1)
	assert.Error(t, err)
}

func TestSetDeckIndex_WrongPlayer_ReturnsError(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	err := p.SetDeckIndex(uuid.NewV4(), 0)
	assert.Error(t, err)
}

// --- SetRole ---

func TestSetRole_MatchesRestriction_Succeeds(t *testing.T) {
	playerID := uuid.NewV4()
	roleID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	p.RoleRestriction = roleID
	err := p.SetRole(playerID, roleID)
	assert.NoError(t, err)
	assert.Equal(t, roleID, p.Role)
}

func TestSetRole_NilRestriction_NilRole_Succeeds(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1) // RoleRestriction = uuid.Nil
	err := p.SetRole(playerID, uuid.Nil)
	assert.NoError(t, err)
}

func TestSetRole_DoesNotMatchRestriction_ReturnsError(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	p.RoleRestriction = uuid.NewV4()
	err := p.SetRole(playerID, uuid.NewV4()) // different role
	assert.Error(t, err)
}

func TestSetRole_WrongPlayer_ReturnsError(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	err := p.SetRole(uuid.NewV4(), uuid.Nil)
	assert.Error(t, err)
}

// --- RemovePlayer ---

func TestRemovePlayer_ClearsPlayerOwnedFields(t *testing.T) {
	playerID := uuid.NewV4()
	p := newParticipant(playerID, 1)
	p.Role = uuid.Nil
	p.DeckIndex = 3
	p.UseStamina = true
	p.FromInvite = true
	p.Ready = true

	p.RemovePlayer()

	assert.False(t, p.HasPlayer())
	assert.Equal(t, uuid.Nil, p.UserID)
	assert.Equal(t, uuid.Nil, p.PlayerID)
	assert.Equal(t, uuid.Nil, p.Role)
	assert.Equal(t, 0, p.DeckIndex)
	assert.False(t, p.UseStamina)
	assert.False(t, p.FromInvite)
	assert.False(t, p.Ready)
}

func TestRemovePlayer_PreservesSlotConfiguration(t *testing.T) {
	lobbyID := uuid.NewV4()
	restriction := uuid.NewV4()
	p := newParticipant(uuid.NewV4(), 2)
	p.LobbyID = lobbyID
	p.RoleRestriction = restriction
	p.InviteOnly = true
	p.Locked = true
	p.BotControl = true

	p.RemovePlayer()

	assert.Equal(t, lobbyID, p.LobbyID)
	assert.Equal(t, 2, p.PlayerSlot)
	assert.Equal(t, restriction, p.RoleRestriction)
	assert.True(t, p.InviteOnly)
	assert.True(t, p.Locked)
	assert.True(t, p.BotControl)
}

func TestRemovePlayer_SlotIsRejoinable(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	p.RemovePlayer()

	joiner := uuid.NewV4()
	err := p.SetPlayer(uuid.NewV4(), joiner, lobby.ParticipantJoinOptions{
		RoleID:     uuid.Nil,
		SlotIndex:  1,
		DeckIndex:  1,
		UseStamina: true,
	})

	assert.NoError(t, err)
	assert.True(t, p.HasPlayer())
	assert.Equal(t, joiner, p.PlayerID)
}

// --- IsHost / HasPlayer ---

func TestIsHost_SlotZero_ReturnsTrue(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 0)
	assert.True(t, p.IsHost())
}

func TestIsHost_NonZeroSlot_ReturnsFalse(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 2)
	assert.False(t, p.IsHost())
}

func TestHasPlayer_BothSet_ReturnsTrue(t *testing.T) {
	p := newParticipant(uuid.NewV4(), 1)
	assert.True(t, p.HasPlayer())
}

func TestHasPlayer_NilPlayerID_ReturnsFalse(t *testing.T) {
	p := &lobby.Participant{
		UserID:   uuid.NewV4(),
		PlayerID: uuid.Nil,
	}
	assert.False(t, p.HasPlayer())
}

func TestHasPlayer_NilUserID_ReturnsFalse(t *testing.T) {
	p := &lobby.Participant{
		UserID:   uuid.Nil,
		PlayerID: uuid.NewV4(),
	}
	assert.False(t, p.HasPlayer())
}
