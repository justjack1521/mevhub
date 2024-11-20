package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/query"
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

	summary, err := g.app.SubApplications.Game.Translators.Summary.Marshall(result)
	if err != nil {
		return nil, err
	}

	return &protomulti.GetGameResponse{
		GameSummary: summary,
	}, nil

}
