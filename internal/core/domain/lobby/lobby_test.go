package lobby_test

import (
	"context"
	"testing"

	"mevhub/internal/core/domain/lobby"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers

func newInstance(t *testing.T) (*lobby.Instance, uuid.UUID) {
	t.Helper()
	hostPlayer := uuid.NewV4()
	hostUser := uuid.NewV4()
	questID := uuid.NewV4()
	lobbyID := uuid.NewV4()

	factory := lobby.NewInstanceFactory(context.Background(), hostUser, hostPlayer)
	opts := lobby.InstanceFactoryOptions{
		QuestID:            questID,
		PlayerSlots:        4,
		MinimumPlayerLevel: 0,
		SlotRestrictions:   make(map[int]lobby.PlayerSlotRestriction),
	}
	inst, err := factory.Create(lobbyID, "12345678", opts)
	require.NoError(t, err)
	return inst, hostPlayer
}

// --- Instance creation ---

func TestNewInstance_NilID_ReturnsError(t *testing.T) {
	factory := lobby.NewInstanceFactory(context.Background(), uuid.NewV4(), uuid.NewV4())
	opts := lobby.InstanceFactoryOptions{
		QuestID:          uuid.NewV4(),
		PlayerSlots:      4,
		SlotRestrictions: make(map[int]lobby.PlayerSlotRestriction),
	}
	_, err := factory.Create(uuid.Nil, "12345678", opts)
	assert.Error(t, err)
}

func TestNewInstance_EmptyPartyID_ReturnsError(t *testing.T) {
	factory := lobby.NewInstanceFactory(context.Background(), uuid.NewV4(), uuid.NewV4())
	opts := lobby.InstanceFactoryOptions{
		QuestID:          uuid.NewV4(),
		PlayerSlots:      4,
		SlotRestrictions: make(map[int]lobby.PlayerSlotRestriction),
	}
	_, err := factory.Create(uuid.NewV4(), "", opts)
	assert.Error(t, err)
}

func TestNewInstance_NilQuestID_ReturnsError(t *testing.T) {
	factory := lobby.NewInstanceFactory(context.Background(), uuid.NewV4(), uuid.NewV4())
	opts := lobby.InstanceFactoryOptions{
		QuestID:          uuid.Nil,
		PlayerSlots:      4,
		SlotRestrictions: make(map[int]lobby.PlayerSlotRestriction),
	}
	_, err := factory.Create(uuid.NewV4(), "12345678", opts)
	assert.Error(t, err)
}

func TestNewInstance_NegativeMinLevel_ReturnsError(t *testing.T) {
	factory := lobby.NewInstanceFactory(context.Background(), uuid.NewV4(), uuid.NewV4())
	opts := lobby.InstanceFactoryOptions{
		QuestID:            uuid.NewV4(),
		PlayerSlots:        4,
		MinimumPlayerLevel: -1,
		SlotRestrictions:   make(map[int]lobby.PlayerSlotRestriction),
	}
	_, err := factory.Create(uuid.NewV4(), "12345678", opts)
	assert.Error(t, err)
}

func TestNewInstance_ValidInputs_HostIsSet(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	assert.Equal(t, hostPlayer, inst.HostPlayerID)
}

// --- StartLobby / CanStart ---

func TestStartLobby_ByHost_Succeeds(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	err := inst.StartLobby(hostPlayer)
	assert.NoError(t, err)
	assert.True(t, inst.Started)
}

func TestStartLobby_ByNonHost_ReturnsError(t *testing.T) {
	inst, _ := newInstance(t)
	nonHost := uuid.NewV4()
	err := inst.StartLobby(nonHost)
	assert.Error(t, err)
	assert.False(t, inst.Started)
}

func TestCanStart_ByHost_ReturnsNil(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	assert.NoError(t, inst.CanStart(hostPlayer))
}

func TestCanStart_ByNonHost_ReturnsError(t *testing.T) {
	inst, _ := newInstance(t)
	assert.Error(t, inst.CanStart(uuid.NewV4()))
}

// --- CanCancel ---

func TestCanCancel_ByHost_ReturnsNil(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	assert.NoError(t, inst.CanCancel(hostPlayer))
}

func TestCanCancel_ByNonHost_ReturnsError(t *testing.T) {
	inst, _ := newInstance(t)
	assert.Error(t, inst.CanCancel(uuid.NewV4()))
}

// --- NewPlayerParticipant ---

func TestNewPlayerParticipant_ValidSlot_Succeeds(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	hostUser := uuid.NewV4()
	opts := lobby.ParticipantJoinOptions{SlotIndex: 0, DeckIndex: 1, UseStamina: true}
	p, err := inst.NewPlayerParticipant(hostUser, hostPlayer, lobby.PlayerSlotRestriction{}, opts)
	require.NoError(t, err)
	assert.Equal(t, 0, p.PlayerSlot)
}

func TestNewPlayerParticipant_SlotBelowZero_ReturnsError(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	opts := lobby.ParticipantJoinOptions{SlotIndex: -1}
	_, err := inst.NewPlayerParticipant(uuid.NewV4(), hostPlayer, lobby.PlayerSlotRestriction{}, opts)
	assert.Error(t, err)
}

func TestNewPlayerParticipant_SlotAtOrAboveMax_ReturnsError(t *testing.T) {
	inst, _ := newInstance(t) // PlayerSlots = 4
	opts := lobby.ParticipantJoinOptions{SlotIndex: 4}
	_, err := inst.NewPlayerParticipant(uuid.NewV4(), uuid.NewV4(), lobby.PlayerSlotRestriction{}, opts)
	assert.Error(t, err)
}

// --- CanAddParticipant ---

func TestCanAddParticipant_HostAtSlotZero_Succeeds(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	opts := lobby.ParticipantJoinOptions{SlotIndex: 0, UseStamina: true}
	p, err := inst.NewPlayerParticipant(uuid.NewV4(), hostPlayer, lobby.PlayerSlotRestriction{}, opts)
	require.NoError(t, err)
	assert.NoError(t, inst.CanAddParticipant(p))
}

func TestCanAddParticipant_HostAtSlotNonZero_ReturnsError(t *testing.T) {
	inst, hostPlayer := newInstance(t)
	opts := lobby.ParticipantJoinOptions{SlotIndex: 1}
	p, err := inst.NewPlayerParticipant(uuid.NewV4(), hostPlayer, lobby.PlayerSlotRestriction{}, opts)
	require.NoError(t, err)
	assert.Error(t, inst.CanAddParticipant(p))
}
