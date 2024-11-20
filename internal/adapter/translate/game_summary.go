package translate

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type GameSummaryTranslator Translator[game.Summary, *protomulti.ProtoGameSummary]
type PartySummaryTranslator Translator[game.PartySummary, *protomulti.ProtoGamePartySummary]

type gameSummaryTranslator struct {
	party PartySummaryTranslator
}

func NewGameSummaryTranslator() GameSummaryTranslator {
	return gameSummaryTranslator{
		party: NewGamePartySummaryTranslator(),
	}
}

func (f gameSummaryTranslator) Marshall(data game.Summary) (out *protomulti.ProtoGameSummary, err error) {
	var result = &protomulti.ProtoGameSummary{
		SysId:   data.SysID.String(),
		Seed:    int32(data.Seed),
		Parties: make([]*protomulti.ProtoGamePartySummary, len(data.Parties)),
	}

	for index, value := range data.Parties {
		party, err := f.party.Marshall(value)
		if err != nil {
			return nil, err
		}
		result.Parties[index] = party
	}
	return result, nil
}

func (f gameSummaryTranslator) Unmarshall(data *protomulti.ProtoGameSummary) (out game.Summary, err error) {
	var result = game.Summary{
		SysID:   uuid.FromStringOrNil(data.SysId),
		Seed:    int(data.Seed),
		Parties: make([]game.PartySummary, len(data.Parties)),
	}
	for index, value := range data.Parties {
		party, err := f.party.Unmarshall(value)
		if err != nil {
			return game.Summary{}, err
		}
		result.Parties[index] = party
	}
	return result, nil
}

type gamePartySummaryTranslator struct {
	player GamePlayerTranslator
}

func NewGamePartySummaryTranslator() PartySummaryTranslator {
	return &gamePartySummaryTranslator{
		player: NewGamePlayerTranslator(),
	}
}

func (f gamePartySummaryTranslator) Marshall(data game.PartySummary) (out *protomulti.ProtoGamePartySummary, err error) {
	var result = &protomulti.ProtoGamePartySummary{
		SysId:     data.SysID.String(),
		PartyId:   data.PartyID,
		Index:     int32(data.Index),
		PartyName: data.PartyName,
		Players:   make([]*protomulti.ProtoGamePlayer, len(data.Players)),
	}
	for index, value := range data.Players {
		player, err := f.player.Marshall(value)
		if err != nil {
			return nil, err
		}
		result.Players[index] = player
	}
	return result, nil
}

func (f gamePartySummaryTranslator) Unmarshall(data *protomulti.ProtoGamePartySummary) (out game.PartySummary, err error) {
	var result = game.PartySummary{
		SysID:     uuid.FromStringOrNil(data.SysId),
		PartyID:   data.PartyId,
		Index:     int(data.Index),
		PartyName: data.PartyName,
		Players:   make([]game.Player, len(data.Players)),
	}
	for index, value := range data.Players {
		player, err := f.player.Unmarshall(value)
		if err != nil {
			return game.PartySummary{}, err
		}
		result.Players[index] = player
	}
	return result, nil
}
