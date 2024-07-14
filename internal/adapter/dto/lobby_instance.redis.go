package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
	"time"
)

type LobbyInstanceRedis struct {
	SysID              string `redis:"SysID"`
	QuestID            string `redis:"QuestID"`
	HostID             string `redis:"HostID"`
	PartyID            string `redis:"PartyID"`
	MinimumPlayerLevel int    `redis:"MinimumPlayerLevel"`
	Started            bool   `redis:"Started"`
	PlayerSlotCount    int    `redis:"PlayerSlotCount"`
	RegisteredAt       int64  `redis:"RegisteredAt"`
}

func (x *LobbyInstanceRedis) ToEntity() *lobby.Instance {
	return &lobby.Instance{
		SysID:              uuid.FromStringOrNil(x.SysID),
		QuestID:            uuid.FromStringOrNil(x.QuestID),
		PartyID:            x.PartyID,
		HostID:             uuid.FromStringOrNil(x.HostID),
		MinimumPlayerLevel: x.MinimumPlayerLevel,
		Started:            x.Started,
		PlayerSlotCount:    x.PlayerSlotCount,
		RegisteredAt:       time.Unix(x.RegisteredAt, 0),
	}
}

func (x *LobbyInstanceRedis) ToMapStringInterface() map[string]interface{} {
	var result = map[string]interface{}{
		"SysID":              x.SysID,
		"QuestID":            x.QuestID,
		"HostID":             x.HostID,
		"PartyID":            x.PartyID,
		"MinimumPlayerLevel": x.MinimumPlayerLevel,
		"Started":            x.Started,
		"PlayerSlotCount":    x.PlayerSlotCount,
		"RegisteredAt":       x.RegisteredAt,
	}
	return result
}
