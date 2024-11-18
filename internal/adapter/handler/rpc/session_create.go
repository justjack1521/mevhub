package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) SessionCreate(context context.Context, request *protomulti.SessionCreateRequest) (*protomulti.SessionCreateResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.SessionCreate.Handle(g.NewCommandContext(context), command.NewSessionCreateCommand()); err != nil {
		return nil, err
	}

	return &protomulti.SessionCreateResponse{}, nil

}
