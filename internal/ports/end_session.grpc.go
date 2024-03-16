package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) EndSession(context context.Context, request *protomulti.EndSessionRequest) (*protomulti.EndSessionResponse, error) {
	ctx, err := g.NewContext(context)
	if err != nil {
		return nil, err
	}
	return g.internal.EndSession(ctx, request)
}

func (g *MultiGrpcServerImplementation) EndSession(ctx GrpcContext, request *protomulti.EndSessionRequest) (*protomulti.EndSessionResponse, error) {

	player, err := uuid.FromString(request.PlayerId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewEndSessionCommand(ctx.ClientID, player)

	if err := g.app.SubApplications.Lobby.Commands.EndSession.Handle(command.NewContext(ctx.Context, ctx.ClientID), cmd); err != nil {
		return nil, err
	}

	return &protomulti.EndSessionResponse{}, nil

}
