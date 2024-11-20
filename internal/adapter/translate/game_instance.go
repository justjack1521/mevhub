package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

type GameInstanceTranslator Translator[*game.Instance, *protomulti.ProtoGameInstance]

type gameInstanceTranslator struct {
}

func NewGameInstanceTranslator() GameInstanceTranslator {
	return gameInstanceTranslator{}
}

func (f gameInstanceTranslator) Marshall(data *game.Instance) (out *protomulti.ProtoGameInstance, err error) {
	var result = &protomulti.ProtoGameInstance{
		SysId:     data.SysID.String(),
		LobbyIds:  make([]string, len(data.LobbyIDs)),
		Seed:      int32(data.Seed),
		State:     int32(data.State),
		StartedAt: data.StartedAt.Unix(),
		Options: &protomulti.ProtoGameInstanceOptions{
			MinimumPlayerLevel: int32(data.Options.MinimumPlayerLevel),
			MaxRunTime:         int64(data.Options.MaxRunTime),
			PlayerTurnDuration: int64(data.Options.PlayerTurnDuration),
			MaxPlayerCount:     int32(data.Options.MaxPlayerCount),
		},
		RegisteredAt: data.RegisteredAt.Unix(),
	}
	for index, value := range data.LobbyIDs {
		result.LobbyIds[index] = value.String()
	}
	return result, nil
}

func (f gameInstanceTranslator) Unmarshall(data *protomulti.ProtoGameInstance) (out *game.Instance, err error) {
	var result = &game.Instance{
		SysID:     uuid.FromStringOrNil(data.SysId),
		LobbyIDs:  make([]uuid.UUID, len(data.LobbyIds)),
		Seed:      int(data.Seed),
		State:     game.InstanceState(data.State),
		StartedAt: time.Unix(data.StartedAt, 0),
		Options: &game.InstanceOptions{
			MinimumPlayerLevel: int(data.Options.MinimumPlayerLevel),
			MaxRunTime:         time.Duration(data.Options.MaxRunTime),
			PlayerTurnDuration: time.Duration(data.Options.PlayerTurnDuration),
			MaxPlayerCount:     int(data.Options.MaxPlayerCount),
		},
		RegisteredAt: time.Unix(data.RegisteredAt, 0),
	}
	for index, value := range data.LobbyIds {
		id, err := uuid.FromString(value)
		if err != nil {
			return nil, err
		}
		result.LobbyIDs[index] = id
	}
	return result, nil
}
