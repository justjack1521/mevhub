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
	RegisteredAt       time.Time
	PlayerSlots        []PlayerSlot
}

type PlayerSlot struct {
	Index           int
	RoleRestriction uuid.UUID
	InviteOnly      bool
	BotControl      bool
	Locked          bool
}

func NewInstance(id uuid.UUID, party string) (*Instance, error) {

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

func (x *Instance) NewPlayerSlot(slot int, restriction PlayerSlotRestriction) (PlayerSlot, error) {
	if slot < 0 || slot >= len(x.PlayerSlots) {
		return PlayerSlot{}, ErrFailedCreateNewParticipant(ErrInvalidPlayerSlot(slot, len(x.PlayerSlots)))
	}
	return PlayerSlot{
		RoleRestriction: restriction.RoleRestriction,
		InviteOnly:      restriction.FromInvite,
		Locked:          restriction.Locked,
		Index:           slot,
	}, nil
}

func (x *Instance) NewPlayerParticipant(client, player uuid.UUID, slot int) (*Participant, error) {
	if slot < 0 || slot >= len(x.PlayerSlots) {
		return nil, ErrFailedCreateNewParticipant(ErrInvalidPlayerSlot(slot, len(x.PlayerSlots)))
	}
	return &Participant{
		ClientID:   client,
		PlayerID:   player,
		LobbyID:    x.SysID,
		PlayerSlot: slot,
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
	if p.ClientID == x.HostID && p.PlayerSlot > 0 {
		return ErrFailedAddParticipant(ErrHostCannotJoinOwnLobby)
	}
	for _, value := range x.PlayerSlots {
		if value.InviteOnly && p.FromInvite == false {
			return ErrFailedAddParticipant(ErrPlayerSlotInviteOnly(value.Index))
		}
	}
	return nil
}

var (
	ErrClientNotLobbyHost = func(client uuid.UUID, id uuid.UUID) error {
		return fmt.Errorf("client %s is not host of lobby %s", client, id)
	}
)

func (x *Instance) CanCancel(client uuid.UUID) error {
	if x.HostID != client {
		return ErrClientNotLobbyHost(client, x.HostID)
	}
	return nil
}
