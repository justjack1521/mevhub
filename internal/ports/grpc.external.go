package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevrpc"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app"
	"mevhub/internal/app/command"
	"mevhub/internal/domain/game"
)

type GrpcContext struct {
	context.Context
}

type MultiGrpcServer struct {
	app *app.CoreApplication
}

func (g MultiGrpcServer) EnqueueAction(ctx context.Context, request *protomulti.GameEnqueueActionRequest) (*protomulti.GameEnqueueActionResponse, error) {

	var cmd = command.NewEnqueueActionCommand(game.PlayerActionType(request.Action), int(request.Target), int(request.SlotIndex), uuid.FromStringOrNil(request.ElementId))

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

func NewMultiGrpcServer(application *app.CoreApplication) MultiGrpcServer {
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
