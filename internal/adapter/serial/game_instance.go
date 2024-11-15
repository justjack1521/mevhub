package serial

import (
	"encoding/json"
	"errors"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"
)

var (
	ErrGameInstanceIsNil = errors.New("nil game instance data passed to serialiser")
)

type GameInstanceSerialiser interface {
	Marshall(data *game.Instance) ([]byte, error)
	Unmarshall(data []byte) (*game.Instance, error)
}

type gameInstanceJSONSerialiser struct {
	translator translate.GameInstanceTranslator
}

func NewGameInstanceJSONSerialiser() GameInstanceSerialiser {
	return gameInstanceJSONSerialiser{translator: translate.NewGameInstanceTranslator()}
}

func (s gameInstanceJSONSerialiser) Marshall(data *game.Instance) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s gameInstanceJSONSerialiser) Unmarshall(data []byte) (*game.Instance, error) {
	if len(data) == 0 {
		return nil, ErrGameInstanceIsNil
	}
	result := &protomulti.ProtoGameInstance{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return s.translator.Unmarshall(result)
}
