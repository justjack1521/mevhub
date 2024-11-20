package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type GameParticipantRedis struct {
	UserID     string `redis:"UserID"`
	PlayerID   string `redis:"PlayerID"`
	PlayerSlot int    `redis:"PlayerSlot"`
	DeckIndex  int    `redis:"DeckIndex"`
	BotControl bool   `redis:"BotControl"`
}

func (x *GameParticipantRedis) ToEntity() *game.Participant {
	return &game.Participant{
		UserID:     uuid.FromStringOrNil(x.UserID),
		PlayerID:   uuid.FromStringOrNil(x.PlayerID),
		PlayerSlot: x.PlayerSlot,
		DeckIndex:  x.DeckIndex,
		BotControl: x.BotControl,
	}
}

func (x *GameParticipantRedis) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"UserID":     x.UserID,
		"PlayerID":   x.PlayerID,
		"PlayerSlot": x.PlayerSlot,
		"DeckIndex":  x.DeckIndex,
		"BotControl": x.BotControl,
	}
}
