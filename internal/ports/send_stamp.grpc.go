package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) SendStamp(ctx context.Context, request *protomulti.SendStampRequest) (*protomulti.SendStampResponse, error) {
	c, err := g.NewContext(ctx)
	if err != nil {
		return nil, err
	}
	return g.internal.SendStamp(c, request)
}

func (g *MultiGrpcServerImplementation) SendStamp(ctx GrpcContext, request *protomulti.SendStampRequest) (*protomulti.SendStampResponse, error) {

	id, err := uuid.FromString(request.StampId)
	if err != nil {
		return nil, err
	}

	var c = command.NewContext(ctx.Context, ctx.ClientID)
	var cmd = command.NewSendStampCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.SendStamp.Handle(c, cmd); err != nil {
		return nil, err
	}

	return &protomulti.SendStampResponse{}, nil

}
