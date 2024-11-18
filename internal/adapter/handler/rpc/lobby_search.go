package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/query"
	"mevhub/internal/core/domain/lobby"
)

func (g MultiGrpcServer) LobbySearch(ctx context.Context, request *protomulti.LobbySearchRequest) (*protomulti.LobbySearchResponse, error) {
	var levels = make([]int, len(request.Levels))
	for index, level := range request.Levels {
		levels[index] = int(level)
	}

	var categories = make([]uuid.UUID, len(request.Categories))
	for index, category := range request.Categories {
		categories[index] = uuid.FromStringOrNil(category)
	}

	var qry = lobby.SearchQuery{
		ModeIdentifier:     request.ModeIdentifier,
		MinimumPlayerLevel: 0,
		Levels:             levels,
		Categories:         categories,
	}

	results, err := g.app.SubApplications.Lobby.Queries.SearchLobby.Handle(g.NewCommandContext(ctx), query.NewSearchLobbyQuery(qry, request.PartyId))
	if err != nil {
		return nil, err
	}

	var lobbies = make([]*protomulti.ProtoLobbySummary, len(results))

	for index, result := range results {
		summary, err := g.app.SubApplications.Lobby.Translators.LobbySummary.Marshall(result)
		if err != nil {
			return nil, err
		}
		lobbies[index] = summary
	}

	return &protomulti.LobbySearchResponse{Lobbies: lobbies}, nil
}
