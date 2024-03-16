package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) CreateSession(context context.Context, request *protomulti.CreateSessionRequest) (*protomulti.CreateSessionResponse, error) {

	ctx, err := g.NewContext(context)
	if err != nil {
		return nil, err
	}
	return g.internal.CreateSession(ctx, request)

}

func (g *MultiGrpcServerImplementation) CreateSession(ctx GrpcContext, request *protomulti.CreateSessionRequest) (*protomulti.CreateSessionResponse, error) {

	player, err := uuid.FromString(request.PlayerId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewCreateSessionCommand(ctx.ClientID, player)

	if err := g.app.SubApplications.Lobby.Commands.CreateSession.Handle(&command.Context{Context: ctx.Context}, cmd); err != nil {
		return nil, err
	}

	return &protomulti.CreateSessionResponse{}, nil

}
