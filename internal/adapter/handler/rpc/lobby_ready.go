package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) LobbyReady(ctx context.Context, request *protomulti.LobbyReadyRequest) (*protomulti.LobbyReadyResponse, error) {

	var cmd = command.NewLobbyReadyCommand()

	if err := g.app.SubApplications.Lobby.Commands.LobbyReady.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.LobbyReadyResponse{}, nil

}
