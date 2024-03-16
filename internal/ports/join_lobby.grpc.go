package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) JoinLobby(context context.Context, request *protomulti.JoinLobbyRequest) (*protomulti.JoinLobbyResponse, error) {
	ctx, err := g.NewContext(context)
	if err != nil {
		return nil, err
	}
	return g.internal.JoinLobby(ctx, request)
}

func (g *MultiGrpcServerImplementation) JoinLobby(ctx GrpcContext, request *protomulti.JoinLobbyRequest) (*protomulti.JoinLobbyResponse, error) {

	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewJoinLobbyCommand(id, int(request.DeckIndex), int(request.SlotIndex), request.UseStamina, request.FromInvite)

	if err := g.app.SubApplications.Lobby.Commands.JoinLobby.Handle(command.NewContext(ctx.Context, ctx.ClientID), cmd); err != nil {
		return nil, err
	}

	return &protomulti.JoinLobbyResponse{}, nil

}
