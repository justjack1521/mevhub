package serial

import (
	"encoding/json"
	"errors"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/domain/lobby"
)

var (
	ErrLobbyPlayerSummaryIsNil = errors.New("nil lobby player summary data passed to serialiser")
)

type LobbyPlayerSummarySerialiser interface {
	Marshall(data lobby.PlayerSummary) ([]byte, error)
	Unmarshall(data []byte) (lobby.PlayerSummary, error)
}

type lobbyPlayerSummaryJSONSerialiser struct {
	translator translate.LobbyPlayerSummaryTranslator
}

func NewLobbyPlayerSummaryJSONSerialiser() LobbyPlayerSummarySerialiser {
	return lobbyPlayerSummaryJSONSerialiser{translator: translate.NewLobbyPlayerSummaryTranslator()}
}

func (s lobbyPlayerSummaryJSONSerialiser) Marshall(data lobby.PlayerSummary) ([]byte, error) {
	p, err := s.translator.Marshall(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (s lobbyPlayerSummaryJSONSerialiser) Unmarshall(data []byte) (lobby.PlayerSummary, error) {
	if len(data) == 0 {
		return lobby.PlayerSummary{}, ErrLobbyPlayerSummaryIsNil
	}
	result := &protomulti.ProtoLobbyPlayer{}
	if err := json.Unmarshal(data, result); err != nil {
		return lobby.PlayerSummary{}, err
	}
	return s.translator.Unmarshall(result)
}
