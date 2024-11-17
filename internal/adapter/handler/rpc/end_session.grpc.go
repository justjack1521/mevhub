package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) EndSession(ctx context.Context, request *protomulti.EndSessionRequest) (*protomulti.EndSessionResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.SessionEnd.Handle(g.NewCommandContext(ctx), command.NewSessionEndCommand()); err != nil {
		return nil, err
	}

	return &protomulti.EndSessionResponse{}, nil
}
