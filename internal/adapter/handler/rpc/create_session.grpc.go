package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) CreateSession(context context.Context, request *protomulti.CreateSessionRequest) (*protomulti.CreateSessionResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.SessionCreate.Handle(g.NewCommandContext(context), command.NewSessionCreateCommand()); err != nil {
		return nil, err
	}

	return &protomulti.CreateSessionResponse{}, nil

}
