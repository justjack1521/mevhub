package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/session"
)

type SessionInstanceRedis struct {
	UserID    string `redis:"UserID"`
	PlayerID  string `redis:"PlayerID"`
	DeckIndex int    `redis:"DeckIndex"`
	LobbyID   string `redis:"LobbyID"`
	GameID    string `redis:"GameID"`
	PartySlot int    `redis:"PartySlot"`
}

func (x *SessionInstanceRedis) ToEntity() *session.Instance {
	return &session.Instance{
		UserID:    uuid.FromStringOrNil(x.UserID),
		PlayerID:  uuid.FromStringOrNil(x.PlayerID),
		DeckIndex: x.DeckIndex,
		LobbyID:   uuid.FromStringOrNil(x.LobbyID),
		GameID:    uuid.FromStringOrNil(x.GameID),
		PartySlot: x.PartySlot,
	}
}

func (x *SessionInstanceRedis) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"UserID":    x.UserID,
		"PlayerID":  x.PlayerID,
		"DeckIndex": x.DeckIndex,
		"LobbyID":   x.LobbyID,
		"GameID":    x.GameID,
		"PartySlot": x.PartySlot,
	}
}
