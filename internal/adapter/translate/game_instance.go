package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
	"time"
)

type GameInstanceTranslator Translator[*game.Instance, *protomulti.ProtoGameInstance]

type gameInstanceTranslator struct {
}

func NewGameInstanceTranslator() GameInstanceTranslator {
	return gameInstanceTranslator{}
}

func (f gameInstanceTranslator) Marshall(data *game.Instance) (out *protomulti.ProtoGameInstance, err error) {
	return &protomulti.ProtoGameInstance{
		SysId:     data.SysID.String(),
		PartyId:   data.PartyID,
		Seed:      int32(data.Seed),
		State:     int32(data.State),
		StartedAt: data.StartedAt.Unix(),
		Options: &protomulti.ProtoGameInstanceOptions{
			MinimumPlayerLevel: int32(data.Options.MinimumPlayerLevel),
			MaxRunTime:         int64(data.Options.MaxRunTime),
			PlayerTurnDuration: int64(data.Options.PlayerTurnDuration),
		},
		RegisteredAt: data.RegisteredAt.Unix(),
	}, nil
}

func (f gameInstanceTranslator) Unmarshall(data *protomulti.ProtoGameInstance) (out *game.Instance, err error) {
	return &game.Instance{
		SysID:     uuid.FromStringOrNil(data.SysId),
		PartyID:   data.PartyId,
		Seed:      int(data.Seed),
		State:     game.InstanceState(data.State),
		StartedAt: time.Unix(data.StartedAt, 0),
		Options: &game.InstanceOptions{
			MinimumPlayerLevel: int(data.Options.MinimumPlayerLevel),
			MaxRunTime:         time.Duration(data.Options.MaxRunTime),
			PlayerTurnDuration: time.Duration(data.Options.PlayerTurnDuration),
		},
		RegisteredAt: time.Unix(data.RegisteredAt, 0),
	}, nil
}
