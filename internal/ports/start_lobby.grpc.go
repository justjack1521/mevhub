package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) StartLobby(ctx context.Context, request *protomulti.StartLobbyRequest) (*protomulti.StartLobbyResponse, error) {

	var cmd = command.NewStartLobbyCommand()

	if err := g.app.SubApplications.Lobby.Commands.StartLobby.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.StartLobbyResponse{}, nil

}
