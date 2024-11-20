package serial

import (
	"encoding/json"
	"errors"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"
)

var (
	ErrGamePlayerParticipantIsNil = errors.New("nil game player participant data passed to serialiser")
)

type GamePlayerSerialiser interface {
	Marshall(data game.Player) ([]byte, error)
	Unmarshall(data []byte) (game.Player, error)
}

type gamePlayerJSONSerialiser struct {
	translator translate.GamePlayerTranslator
}

func NewGamePlayerJSONSerialiser() GamePlayerSerialiser {
	return &gamePlayerJSONSerialiser{translator: translate.NewGamePlayerTranslator()}
}

func (s gamePlayerJSONSerialiser) Marshall(data game.Player) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s gamePlayerJSONSerialiser) Unmarshall(data []byte) (game.Player, error) {
	if len(data) == 0 {
		return game.Player{}, ErrGamePlayerParticipantIsNil
	}
	result := &protomulti.ProtoGamePlayer{}
	if err := json.Unmarshal(data, result); err != nil {
		return game.Player{}, err
	}
	return s.translator.Unmarshall(result)
}
