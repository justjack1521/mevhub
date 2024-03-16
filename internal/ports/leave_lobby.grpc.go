package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
)

func (g MultiGrpcServer) LeaveLobby(ctx context.Context, request *protomulti.LeaveLobbyRequest) (*protomulti.LeaveLobbyResponse, error) {
	c, err := g.NewContext(ctx)
	if err != nil {
		return nil, err
	}
	return g.internal.LeaveLobby(c, request)
}

func (g *MultiGrpcServerImplementation) LeaveLobby(ctx GrpcContext, request *protomulti.LeaveLobbyRequest) (*protomulti.LeaveLobbyResponse, error) {
	panic(nil)
}
