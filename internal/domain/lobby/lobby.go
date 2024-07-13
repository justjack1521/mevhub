package lobby

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	ErrFailedCreateNewInstance = func(err error) error {
		return fmt.Errorf("failed to create new lobby instance: %w", err)
	}
	ErrHostIDNil     = errors.New("host id is nil")
	ErrInstanceIDNil = errors.New("instance id is nil")
	ErrPartyIDEmpty  = errors.New("party id is empty")
)

type Instance struct {
	SysID              uuid.UUID
	QuestID            uuid.UUID
	HostID             uuid.UUID
	PartyID            string
	MinimumPlayerLevel int
	Started            bool
	PlayerSlotCount    int
	RegisteredAt       time.Time
}

func NewInstance(id uuid.UUID, party string, options InstanceFactoryOptions) (*Instance, error) {

	if id == uuid.Nil {
		return nil, ErrFailedCreateNewInstance(ErrInstanceIDNil)
	}

	if party == "" {
		return nil, ErrFailedCreateNewInstance(ErrPartyIDEmpty)
	}

	return &Instance{
		SysID:              id,
		MinimumPlayerLevel: 0,
		PartyID:            party,
		RegisteredAt:       time.Now().UTC(),
		PlayerSlotCount:    options.PlayerSlots,
	}, nil

}

var (
	ErrQuestIDNil      = errors.New("quest id is nil")
	ErrMinLevelInvalid = func(level int) error {
		return fmt.Errorf("min player level is invalid: %d", level)
	}
)

func (x *Instance) SetQuestID(id uuid.UUID) error {
	if uuid.Equal(id, uuid.Nil) {
		return ErrQuestIDNil
	}
	x.QuestID = id
	return nil
}

func (x *Instance) SetMinPlayerLevel(level int) error {
	if level < 0 {
		return ErrMinLevelInvalid(level)
	}
	x.MinimumPlayerLevel = level
	return nil
}

var (
	ErrFailedCreateNewParticipant = func(err error) error {
		return fmt.Errorf("failed to create new participant instance: %w", err)
	}
	ErrInvalidPlayerSlot = func(index int, max int) error {
		return fmt.Errorf("invalid player slot: %d max slots: %d", index, max)
	}
)

func (x *Instance) NewPlayerParticipant(client, player uuid.UUID, restriction PlayerSlotRestriction, options ParticipantJoinOptions) (*Participant, error) {
	if options.SlotIndex < 0 || options.SlotIndex >= x.PlayerSlotCount {
		return nil, ErrFailedCreateNewParticipant(ErrInvalidPlayerSlot(options.SlotIndex, x.PlayerSlotCount))
	}
	return &Participant{
		UserID:          client,
		PlayerID:        player,
		LobbyID:         x.SysID,
		Role:            options.RoleID,
		RoleRestriction: restriction.RoleRestriction,
		InviteOnly:      restriction.InviteOnly,
		Locked:          restriction.Locked,
		PlayerSlot:      options.SlotIndex,
		DeckIndex:       options.DeckIndex,
		UseStamina:      options.UseStamina,
		FromInvite:      options.FromInvite,
		BotControl:      restriction.BotControl,
	}, nil
}

var (
	ErrFailedAddParticipant = func(err error) error {
		return fmt.Errorf("failed to add lobby participant: %w", err)
	}
	ErrParticipantExistsInLobby = func(index int) error {
		return fmt.Errorf("participants already added to slot index %d", index)
	}
	ErrPlayerSlotUnavailable = func(index int) error {
		return fmt.Errorf("player slot %d is unavailable", index)
	}
	ErrHostCannotJoinOwnLobby = errors.New("host cannot join own lobby")
	ErrPlayerSlotInviteOnly   = func(index int) error {
		return fmt.Errorf("player slot %d is invite only", index)
	}
)

func (x *Instance) CanAddParticipant(p *Participant) error {
	if p.UserID == x.HostID && p.PlayerSlot > 0 {
		return ErrFailedAddParticipant(ErrHostCannotJoinOwnLobby)
	}
	return nil
}

var (
	ErrClientNotLobbyHost = func(user uuid.UUID, id uuid.UUID) error {
		return fmt.Errorf("client %s is not host of lobby %s", user, id)
	}
)

func (x *Instance) CanCancel(user uuid.UUID) error {
	if x.HostID != user {
		return ErrClientNotLobbyHost(user, x.HostID)
	}
	return nil
}
