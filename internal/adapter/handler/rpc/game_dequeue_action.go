package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) DequeueAction(ctx context.Context, request *protomulti.GameDequeueActionRequest) (*protomulti.GameDequeueActionResponse, error) {

	var cmd = command.NewDequeueActionCommand()

	if err := g.app.SubApplications.Game.Commands.DequeueAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameDequeueActionResponse{}, nil

}
