package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) PlayerDeath(ctx context.Context, request *protomulti.GamePlayerDeathRequest) (*protomulti.GamePlayerDeathResponse, error) {

	var cmd = command.NewPlayerDeathCommand()

	if err := g.app.SubApplications.Game.Commands.PlayerDeath.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GamePlayerDeathResponse{}, nil

}
