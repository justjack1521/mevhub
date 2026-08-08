package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/session"
)

type SessionInstanceRedis struct {
	UserID              string `redis:"UserID"`
	PlayerID            string `redis:"PlayerID"`
	DeckIndex           int    `redis:"DeckIndex"`
	LobbyID             string `redis:"LobbyID"`
	GameID              string `redis:"GameID"`
	PartySlot           int    `redis:"PartySlot"`
	DisconnectSessionID string `redis:"DisconnectSessionID"`
	LastConnEventAt     int64  `redis:"LastConnEventAt"`
	DisconnectedAt      int64  `redis:"DisconnectedAt"`
	CurrentSessionID    string `redis:"CurrentSessionID"`
}

func (x *SessionInstanceRedis) ToEntity() *session.Instance {
	return &session.Instance{
		UserID:              uuid.FromStringOrNil(x.UserID),
		PlayerID:            uuid.FromStringOrNil(x.PlayerID),
		DeckIndex:           x.DeckIndex,
		LobbyID:             uuid.FromStringOrNil(x.LobbyID),
		GameID:              uuid.FromStringOrNil(x.GameID),
		PartySlot:           x.PartySlot,
		DisconnectSessionID: uuid.FromStringOrNil(x.DisconnectSessionID),
		LastConnEventAt:     x.LastConnEventAt,
		DisconnectedAt:      x.DisconnectedAt,
		CurrentSessionID:    uuid.FromStringOrNil(x.CurrentSessionID),
	}
}

func (x *SessionInstanceRedis) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"UserID":              x.UserID,
		"PlayerID":            x.PlayerID,
		"DeckIndex":           x.DeckIndex,
		"LobbyID":             x.LobbyID,
		"GameID":              x.GameID,
		"PartySlot":           x.PartySlot,
		"DisconnectSessionID": x.DisconnectSessionID,
		"LastConnEventAt":     x.LastConnEventAt,
		"DisconnectedAt":      x.DisconnectedAt,
		"CurrentSessionID":    x.CurrentSessionID,
	}
}

// ConnectionStateMap holds only the fields the connect/disconnect handlers
// mutate, for narrow HSET writes that cannot clobber concurrent lobby/game
// field updates.
func (x *SessionInstanceRedis) ConnectionStateMap() map[string]interface{} {
	return map[string]interface{}{
		"DisconnectSessionID": x.DisconnectSessionID,
		"LastConnEventAt":     x.LastConnEventAt,
		"DisconnectedAt":      x.DisconnectedAt,
		"CurrentSessionID":    x.CurrentSessionID,
	}
}
