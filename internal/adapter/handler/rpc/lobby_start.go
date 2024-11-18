package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) LobbyStart(ctx context.Context, request *protomulti.LobbyStartRequest) (*protomulti.LobbyStartResponse, error) {

	var cmd = command.NewLobbyStartCommand()

	if err := g.app.SubApplications.Lobby.Commands.LobbyStart.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.LobbyStartResponse{}, nil

}
