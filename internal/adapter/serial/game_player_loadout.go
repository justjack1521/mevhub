package serial

import (
	"encoding/json"
	"errors"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"

	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
)

var (
	ErrGamePlayerLoadoutIsNil = errors.New("nil game player loadout data passed to serialiser")
)

type GamePlayerLoadoutSerialiser interface {
	Marshall(data game.PlayerLoadout) ([]byte, error)
	Unmarshall(data []byte) (game.PlayerLoadout, error)
}

type gamePlayerLoadoutJSONSerialiser struct {
	translator translate.GamePlayerLoadoutTranslator
}

func NewGamePlayerLoadoutJSONSerialiser() GamePlayerLoadoutSerialiser {
	return &gamePlayerLoadoutJSONSerialiser{translator: translate.NewGamePlayerLoadoutTranslator()}
}

func (s gamePlayerLoadoutJSONSerialiser) Marshall(data game.PlayerLoadout) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s gamePlayerLoadoutJSONSerialiser) Unmarshall(data []byte) (game.PlayerLoadout, error) {
	if len(data) == 0 {
		return game.PlayerLoadout{}, ErrGamePlayerLoadoutIsNil
	}
	result := &protoidentity.ProtoPlayerLoadout{}
	if err := json.Unmarshal(data, result); err != nil {
		return game.PlayerLoadout{}, err
	}
	return s.translator.Unmarshall(result)
}
