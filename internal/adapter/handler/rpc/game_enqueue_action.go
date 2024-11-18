package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/domain/game"
)

func (g MultiGrpcServer) EnqueueAction(ctx context.Context, request *protomulti.GameEnqueueActionRequest) (*protomulti.GameEnqueueActionResponse, error) {

	var cmd = command.NewEnqueueActionCommand(game.PlayerActionType(request.Action.Action), int(request.Action.Target), int(request.Action.SlotIndex), uuid.FromStringOrNil(request.Action.ElementId))

	if err := g.app.SubApplications.Game.Commands.EnqueueAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameEnqueueActionResponse{}, nil

}
