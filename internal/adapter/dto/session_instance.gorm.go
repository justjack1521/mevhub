package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/session"
)

type SessionInstanceGorm struct {
	UserID    uuid.UUID `gorm:"primaryKey;column:client_id"`
	PlayerID  uuid.UUID `gorm:"column:player_id"`
	DeckIndex int       `gorm:"column:deck_index"`
}

func (SessionInstanceGorm) TableName() string {
	return "lobby.session"
}

func (x *SessionInstanceGorm) ToEntity() *session.Instance {
	return &session.Instance{
		UserID:    x.UserID,
		PlayerID:  x.PlayerID,
		DeckIndex: x.DeckIndex,
	}
}
