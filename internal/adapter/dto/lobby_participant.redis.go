package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbyParticipantRedis struct {
	ClientID        string `redis:"UserID"`
	PlayerID        string `redis:"PlayerID"`
	LobbyID         string `redis:"LobbyID"`
	Role            string `redis:"Role"`
	RoleRestriction string `redis:"RoleRestriction"`
	Locked          bool   `redis:"Locked"`
	InviteOnly      bool   `redis:"InviteOnly"`
	PlayerSlot      int    `redis:"PlayerSlot"`
	DeckIndex       int    `redis:"DeckIndex"`
	UseStamina      bool   `redis:"UseStamina"`
	FromInvite      bool   `redis:"FromInvite"`
	Ready           bool   `redis:"Ready"`
}

func (x *LobbyParticipantRedis) ToEntity() *lobby.Participant {
	return &lobby.Participant{
		UserID:          uuid.FromStringOrNil(x.ClientID),
		PlayerID:        uuid.FromStringOrNil(x.PlayerID),
		LobbyID:         uuid.FromStringOrNil(x.LobbyID),
		RoleRestriction: uuid.FromStringOrNil(x.RoleRestriction),
		Locked:          x.Locked,
		InviteOnly:      x.InviteOnly,
		Role:            uuid.FromStringOrNil(x.Role),
		PlayerSlot:      x.PlayerSlot,
		DeckIndex:       x.DeckIndex,
		UseStamina:      x.UseStamina,
		FromInvite:      x.FromInvite,
		Ready:           x.Ready,
	}
}

func (x *LobbyParticipantRedis) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"UserID":          x.ClientID,
		"PlayerID":        x.PlayerID,
		"LobbyID":         x.LobbyID,
		"RoleRestriction": x.RoleRestriction,
		"Locked":          x.Locked,
		"InviteOnly":      x.InviteOnly,
		"Role":            x.Role,
		"PlayerSlot":      x.PlayerSlot,
		"DeckIndex":       x.DeckIndex,
		"UseStamina":      x.UseStamina,
		"FromInvite":      x.FromInvite,
		"Ready":           x.Ready,
	}
}
