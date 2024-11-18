package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/query"
)

func (g MultiGrpcServer) GetGame(ctx context.Context, request *protomulti.GetGameRequest) (*protomulti.GetGameResponse, error) {

	var qry = query.NewGameSummaryQuery()

	result, err := g.app.SubApplications.Game.Queries.GameSummary.Handle(g.NewCommandContext(ctx), qry)
	if err != nil {
		return nil, err
	}

	var participants = make([]*protomulti.ProtoGameParticipant, len(result.Participants))

	for index, value := range result.Participants {
		participant, err := g.app.SubApplications.Game.Translators.PlayerParticipant.Marshall(value)
		if err != nil {
			return nil, err
		}
		participants[index] = participant
	}

	return &protomulti.GetGameResponse{
		GameData: &protomulti.ProtoGameInstance{
			SysId:   result.SysID.String(),
			PartyId: result.PartyID,
			Seed:    int32(result.Seed),
		},
		Participants: participants,
	}, nil

}
