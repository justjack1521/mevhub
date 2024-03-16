package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/server"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app"
)

type GrpcContext struct {
	context.Context
	ClientID uuid.UUID
}

type MultiGrpcServer struct {
	internal *MultiGrpcServerImplementation
	app      *app.Application
}

func NewMultiGrpcServer(application *app.Application) MultiGrpcServer {
	return MultiGrpcServer{app: application, internal: NewMultiGrpcServerImplementation(application)}
}

func (g MultiGrpcServer) NewContext(ctx context.Context) (GrpcContext, error) {
	client, err := server.ExtractUserIDFromContext(ctx)
	if err != nil {
		return GrpcContext{}, err
	}
	return GrpcContext{
		Context:  ctx,
		ClientID: client,
	}, nil
}
