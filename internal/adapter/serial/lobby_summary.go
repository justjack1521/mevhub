package serial

import (
	"encoding/json"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/lobby"
)

type LobbySummarySerialiser interface {
	Marshall(data lobby.Summary) ([]byte, error)
	Unmarshall(data []byte) (lobby.Summary, error)
}

type lobbySummaryJSONSerialiser struct {
	translator translate.LobbySummaryTranslator
}

func NewLobbySummaryJSONSerialiser() LobbySummarySerialiser {
	return &lobbySummaryJSONSerialiser{translator: translate.NewLobbySummaryTranslator()}
}

func (s lobbySummaryJSONSerialiser) Marshall(data lobby.Summary) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s lobbySummaryJSONSerialiser) Unmarshall(data []byte) (lobby.Summary, error) {
	if len(data) == 0 {
		return lobby.Summary{}, ErrLobbyPlayerSummaryIsNil
	}
	result := &protomulti.ProtoLobbySummary{}
	if err := json.Unmarshal(data, result); err != nil {
		return lobby.Summary{}, err
	}
	return s.translator.Unmarshall(result)
}
