package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevrpc"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/domain/game"
)

type GrpcContext struct {
	context.Context
}

type MultiGrpcServer struct {
	app *application.CoreApplication
}

func (g MultiGrpcServer) ReadyPlayer(ctx context.Context, request *protomulti.GameReadyPlayerRequest) (*protomulti.GameReadyPlayerResponse, error) {

	var cmd = command.NewReadyPlayerCommand()

	if err := g.app.SubApplications.Game.Commands.ReadyPlayer.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameReadyPlayerResponse{}, nil

}

func (g MultiGrpcServer) EnqueueAction(ctx context.Context, request *protomulti.GameEnqueueActionRequest) (*protomulti.GameEnqueueActionResponse, error) {

	var cmd = command.NewEnqueueActionCommand(game.PlayerActionType(request.Action.Action), int(request.Action.Target), int(request.Action.SlotIndex), uuid.FromStringOrNil(request.Action.ElementId))

	if err := g.app.SubApplications.Game.Commands.EnqueueAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameEnqueueActionResponse{}, nil

}

func (g MultiGrpcServer) DequeueAction(ctx context.Context, request *protomulti.GameDequeueActionRequest) (*protomulti.GameDequeueActionResponse, error) {

	var cmd = command.NewDequeueActionCommand()

	if err := g.app.SubApplications.Game.Commands.DequeueAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameDequeueActionResponse{}, nil

}

func (g MultiGrpcServer) LockAction(ctx context.Context, request *protomulti.GameLockActionRequest) (*protomulti.GameLockActionResponse, error) {

	var cmd = command.NewLockActionCommand()

	if err := g.app.SubApplications.Game.Commands.LockAction.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameLockActionResponse{}, nil

}

func NewMultiGrpcServer(application *application.CoreApplication) MultiGrpcServer {
	return MultiGrpcServer{app: application}
}

func (g GrpcContext) UserID() uuid.UUID {
	return mevrpc.UserIDFromContext(g.Context)
}

func (g GrpcContext) PlayerID() uuid.UUID {
	return mevrpc.PlayerIDFromContext(g.Context)
}

func (g MultiGrpcServer) NewCommandContext(ctx context.Context) GrpcContext {
	return GrpcContext{Context: ctx}
}
