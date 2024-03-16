package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) ReadyLobby(ctx context.Context, request *protomulti.ReadyLobbyRequest) (*protomulti.ReadyLobbyResponse, error) {
	c, err := g.NewContext(ctx)
	if err != nil {
		return nil, err
	}
	return g.internal.ReadyLobby(c, request)
}

func (g *MultiGrpcServerImplementation) ReadyLobby(ctx GrpcContext, request *protomulti.ReadyLobbyRequest) (*protomulti.ReadyLobbyResponse, error) {

	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewReadyLobbyCommand(id, int(request.DeckIndex))

	if err := g.app.SubApplications.Lobby.Commands.ReadyLobby.Handle(command.NewContext(ctx.Context, ctx.ClientID), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ReadyLobbyResponse{}, nil

}
