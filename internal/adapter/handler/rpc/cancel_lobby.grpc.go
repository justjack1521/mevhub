package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) CancelLobby(ctx context.Context, request *protomulti.CancelLobbyRequest) (*protomulti.CancelLobbyResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.LobbyCancel.Handle(g.NewCommandContext(ctx), command.NewLobbyCancelCommand()); err != nil {
		return nil, err
	}

	return &protomulti.CancelLobbyResponse{}, nil
}
