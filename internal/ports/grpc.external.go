package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevrpc"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app"
)

type GrpcContext struct {
	context.Context
}

type MultiGrpcServer struct {
	app *app.CoreApplication
}

func (g MultiGrpcServer) EnqueueAbility(ctx context.Context, request *protomulti.GameEnqueueAbilityRequest) (*protomulti.GameEnqueueAbilityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g MultiGrpcServer) DequeueAbility(ctx context.Context, request *protomulti.GameDequeueAbilityRequest) (*protomulti.GameDequeueAbilityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g MultiGrpcServer) LockAction(ctx context.Context, request *protomulti.GameLockActionRequest) (*protomulti.GameLockActionResponse, error) {
	//TODO implement me
	panic("implement me")
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
