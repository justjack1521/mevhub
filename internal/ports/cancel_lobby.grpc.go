package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) CancelLobby(ctx context.Context, request *protomulti.CancelLobbyRequest) (*protomulti.CancelLobbyResponse, error) {
	c, err := g.NewContext(ctx)
	if err != nil {
		return nil, err
	}
	return g.internal.CancelLobby(c, request)
}

func (g *MultiGrpcServerImplementation) CancelLobby(ctx GrpcContext, request *protomulti.CancelLobbyRequest) (*protomulti.CancelLobbyResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.CancelLobby.Handle(command.NewContext(ctx.Context, ctx.ClientID), command.NewCancelLobbyCommand()); err != nil {
		return nil, err
	}

	return &protomulti.CancelLobbyResponse{}, nil

}
