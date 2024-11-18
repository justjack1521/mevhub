package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ReadyPlayer(ctx context.Context, request *protomulti.GameReadyPlayerRequest) (*protomulti.GameReadyPlayerResponse, error) {

	var cmd = command.NewReadyPlayerCommand()

	if err := g.app.SubApplications.Game.Commands.ReadyPlayer.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameReadyPlayerResponse{}, nil

}
