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
		SysId:        data.SysID.String(),
		PartyId:      data.PartyID,
		Seed:         data.Seed,
		State:        int32(data.State),
		StartedAt:    data.StartedAt.Unix(),
		RegisteredAt: data.RegisteredAt.Unix(),
	}, nil
}

func (f gameInstanceTranslator) Unmarshall(data *protomulti.ProtoGameInstance) (out *game.Instance, err error) {
	return &game.Instance{
		SysID:        uuid.FromStringOrNil(data.SysId),
		PartyID:      data.PartyId,
		Seed:         data.Seed,
		State:        game.InstanceState(data.State),
		StartedAt:    time.Unix(data.StartedAt, 0),
		RegisteredAt: time.Unix(data.RegisteredAt, 0),
	}, nil
}