package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) EndSession(ctx context.Context, request *protomulti.EndSessionRequest) (*protomulti.EndSessionResponse, error) {

	if err := g.app.SubApplications.Lobby.Commands.EndSession.Handle(g.NewCommandContext(ctx), command.NewEndSessionCommand()); err != nil {
		return nil, err
	}

	return &protomulti.EndSessionResponse{}, nil
}
