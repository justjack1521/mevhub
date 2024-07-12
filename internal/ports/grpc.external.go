package ports

import (
	"context"
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
