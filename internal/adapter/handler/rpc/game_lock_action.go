package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) LockAction(ctx context.Context, request *protomulti.GameLockActionRequest) (*protomulti.GameLockActionResponse, error) {

	var cmd = command.NewLockActionCommand()

	if err := g.app.SubApplications.Game.Commands.LockAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameLockActionResponse{}, nil

}
