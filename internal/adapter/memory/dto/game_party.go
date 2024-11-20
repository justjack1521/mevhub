package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type GamePartyRedis struct {
	SysID     string `redis:"SysID"`
	PartyID   string `redis:"PartyID"`
	Index     int    `redis:"Index"`
	PartyName string `redis:"PartyName"`
}

func (x *GamePartyRedis) ToEntity() *game.Party {
	return &game.Party{
		SysID:     uuid.FromStringOrNil(x.SysID),
		PartyID:   x.PartyID,
		Index:     x.Index,
		PartyName: x.PartyName,
	}
}

func (x *GamePartyRedis) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"SysID":     x.SysID,
		"PartyID":   x.PartyID,
		"Index":     x.Index,
		"PartyName": x.PartyName,
	}
}
