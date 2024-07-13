package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbyParticipantGorm struct {
	ClientID       uuid.UUID `gorm:"column:client_id"`
	PlayerID       uuid.UUID `gorm:"column:player_id"`
	GameInstanceID uuid.UUID `gorm:"column:game_instance"`
	Role           uuid.UUID `gorm:"column:role"`
	PlayerSlot     int       `gorm:"column:player_slot"`
	DeckIndex      int       `gorm:"column:deck_index"`
	UseStamina     bool      `gorm:"column:use_stamina"`
	FromInvite     bool      `gorm:"column:from_invite"`
	Ready          bool      `gorm:"column:ready"`
}

func (LobbyParticipantGorm) TableName() string {
	return "multi.lobby_instance_participant"
}

func (x *LobbyParticipantGorm) ToEntity() *lobby.Participant {
	return &lobby.Participant{
		UserID:     x.ClientID,
		PlayerID:   x.PlayerID,
		LobbyID:    x.GameInstanceID,
		Role:       x.Role,
		PlayerSlot: x.PlayerSlot,
		DeckIndex:  x.DeckIndex,
		UseStamina: x.UseStamina,
		FromInvite: x.FromInvite,
	}
}
