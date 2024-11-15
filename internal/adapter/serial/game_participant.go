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

type GamePlayerParticipantSerialiser interface {
	Marshall(data *game.PlayerParticipant) ([]byte, error)
	Unmarshall(data []byte) (*game.PlayerParticipant, error)
}

type gamePlayerParticipantJSONSerialiser struct {
	translator translate.GamePlayerParticipantTranslator
}

func NewGamePlayerParticipantJSONSerialiser() GamePlayerParticipantSerialiser {
	return &gamePlayerParticipantJSONSerialiser{translator: translate.NewGameParticipantTranslator()}
}

func (s gamePlayerParticipantJSONSerialiser) Marshall(data *game.PlayerParticipant) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s gamePlayerParticipantJSONSerialiser) Unmarshall(data []byte) (*game.PlayerParticipant, error) {
	if len(data) == 0 {
		return nil, ErrGamePlayerParticipantIsNil
	}
	result := &protomulti.ProtoGameParticipant{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return s.translator.Unmarshall(result)
}
