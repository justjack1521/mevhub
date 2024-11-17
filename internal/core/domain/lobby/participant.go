package lobby

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrParticipantNotBelongToPlayer = func(actual, expected uuid.UUID) error {
		return fmt.Errorf("participant %s does not belong to player %s", expected, actual)
	}
	ErrInvalidParticipantDeckIndex = func(index int) error {
		return fmt.Errorf("participant deck index %d must be non-negative", index)
	}
)

type Participant struct {
	UserID          uuid.UUID
	PlayerID        uuid.UUID
	LobbyID         uuid.UUID
	Role            uuid.UUID
	RoleRestriction uuid.UUID
	InviteOnly      bool
	Locked          bool
	PlayerSlot      int
	DeckIndex       int
	UseStamina      bool
	FromInvite      bool
	Ready           bool
	BotControl      bool
}

func (x *Participant) IsHost() bool {
	return x.PlayerSlot == 0
}

func (x *Participant) HasPlayer() bool {
	return x.UserID != uuid.Nil && x.PlayerID != uuid.Nil
}

func (x *Participant) SetPlayer(user, player uuid.UUID, options ParticipantJoinOptions) error {

	x.UserID = user
	x.PlayerID = player

	if err := x.SetRole(player, options.RoleID); err != nil {
		return err
	}

	if err := x.SetUseStamina(player, options.UseStamina); err != nil {
		return err
	}

	if err := x.SetDeckIndex(player, options.DeckIndex); err != nil {
		return err
	}

	return nil

}

func (x *Participant) SetReady(player uuid.UUID, value bool) error {
	if uuid.Equal(player, x.PlayerID) == false {
		return ErrParticipantNotBelongToPlayer(player, x.PlayerID)
	}
	x.Ready = value
	return nil
}

var (
	ErrHostMustUseStamina = errors.New("host must use stamina")
)

func (x *Participant) SetUseStamina(player uuid.UUID, value bool) error {
	if uuid.Equal(player, x.PlayerID) == false {
		return ErrParticipantNotBelongToPlayer(player, x.PlayerID)
	}
	if value == false && x.PlayerSlot == 0 {
		return ErrHostMustUseStamina
	}
	x.UseStamina = value
	return nil
}

var (
	ErrRoleNotAllow = func(attempt uuid.UUID, expected uuid.UUID) error {
		return fmt.Errorf("role %s not allowed, expected %s", attempt, expected)
	}
)

func (x *Participant) SetRole(player uuid.UUID, role uuid.UUID) error {
	if uuid.Equal(player, x.PlayerID) == false {
		return ErrParticipantNotBelongToPlayer(player, x.PlayerID)
	}
	if uuid.Equal(role, x.RoleRestriction) == false {
		return ErrRoleNotAllow(role, x.RoleRestriction)
	}
	x.Role = role
	return nil
}

func (x *Participant) SetDeckIndex(player uuid.UUID, index int) error {
	if uuid.Equal(player, x.PlayerID) == false {
		return ErrParticipantNotBelongToPlayer(player, x.PlayerID)
	}
	if index < 0 {
		return ErrInvalidParticipantDeckIndex(index)
	}
	x.DeckIndex = index
	return nil
}
