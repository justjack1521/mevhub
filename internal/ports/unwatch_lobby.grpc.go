package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
)

func (g MultiGrpcServer) UnwatchLobby(ctx context.Context, request *protomulti.UnwatchLobbyRequest) (*protomulti.UnwatchLobbyResponse, error) {
	c, err := g.NewContext(ctx)
	if err != nil {
		return nil, err
	}
	return g.internal.UnwatchLobby(c, request)
}

func (g *MultiGrpcServerImplementation) UnwatchLobby(ctx GrpcContext, request *protomulti.UnwatchLobbyRequest) (*protomulti.UnwatchLobbyResponse, error) {
	panic(nil)
}
