package dto

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
	"time"
)

type LobbyInstanceRedis struct {
	LobbyID            string `redis:"LobbyID"`
	QuestID            string `redis:"QuestID"`
	PartyID            string `redis:"PartyID"`
	HostID             string `redis:"HostID"`
	MinimumPlayerLevel int    `redis:"MinimumPlayerLevel"`
	Started            bool   `redis:"Started"`
	RegisteredAt       int64  `redis:"RegisteredAt"`
}

func (x *LobbyInstanceRedis) ToEntity() *lobby.Instance {
	return &lobby.Instance{
		SysID:              uuid.FromStringOrNil(x.LobbyID),
		QuestID:            uuid.FromStringOrNil(x.QuestID),
		PartyID:            x.PartyID,
		HostID:             uuid.FromStringOrNil(x.HostID),
		MinimumPlayerLevel: x.MinimumPlayerLevel,
		RegisteredAt:       time.Unix(x.RegisteredAt, 0),
	}
}

func (x *LobbyInstanceRedis) ToMapStringInterface() map[string]interface{} {
	var result = map[string]interface{}{
		"LobbyID":            x.LobbyID,
		"QuestID":            x.QuestID,
		"PartyID":            x.PartyID,
		"HostID":             x.HostID,
		"MinimumPlayerLevel": x.MinimumPlayerLevel,
		"Started":            x.Started,
		"RegisteredAt":       x.RegisteredAt,
	}
	return result
}
