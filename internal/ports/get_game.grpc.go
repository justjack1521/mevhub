package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/query"
)

func (g MultiGrpcServer) GetGame(ctx context.Context, request *protomulti.GetGameRequest) (*protomulti.GetGameResponse, error) {

	id, err := uuid.FromString(request.GameId)
	if err != nil {
		return nil, err
	}

	var qry = query.NewGameSummaryQuery(id)

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
			Seed:    result.Seed,
		},
		Participants: participants,
	}, nil

}
