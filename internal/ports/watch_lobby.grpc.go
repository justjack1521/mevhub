package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) WatchLobby(context context.Context, request *protomulti.WatchLobbyRequest) (*protomulti.WatchLobbyResponse, error) {
	ctx, err := g.NewContext(context)
	if err != nil {
		return nil, err
	}
	return g.internal.WatchLobby(ctx, request)
}

func (g *MultiGrpcServerImplementation) WatchLobby(ctx GrpcContext, request *protomulti.WatchLobbyRequest) (*protomulti.WatchLobbyResponse, error) {

	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewWatchLobbyCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.WatchLobby.Handle(command.NewContext(ctx.Context, ctx.ClientID), cmd); err != nil {
		return nil, err
	}

	return &protomulti.WatchLobbyResponse{}, nil

}
