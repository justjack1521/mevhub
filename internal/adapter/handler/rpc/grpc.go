package rpc

import (
	"context"
	"github.com/justjack1521/mevrpc"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application"
)

type MultiGrpcServer struct {
	app *application.CoreApplication
}

func NewMultiGrpcServer(application *application.CoreApplication) MultiGrpcServer {
	return MultiGrpcServer{app: application}
}

type GrpcContext struct {
	context.Context
}

func (g MultiGrpcServer) NewCommandContext(ctx context.Context) GrpcContext {
	return GrpcContext{Context: ctx}
}

func (g GrpcContext) UserID() uuid.UUID {
	return mevrpc.UserIDFromContext(g.Context)
}

func (g GrpcContext) PlayerID() uuid.UUID {
	return mevrpc.PlayerIDFromContext(g.Context)
}
