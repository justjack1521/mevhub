package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) GameStamp(ctx context.Context, request *protomulti.GameStampRequest) (*protomulti.GameStampResponse, error) {

	id, err := uuid.FromString(request.StampId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewGameStampCommand(id)

	if err := g.app.SubApplications.Game.Commands.GameStamp.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameStampResponse{}, nil

}
