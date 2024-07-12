package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) WatchLobby(ctx context.Context, request *protomulti.WatchLobbyRequest) (*protomulti.WatchLobbyResponse, error) {
	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewWatchLobbyCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.WatchLobby.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.WatchLobbyResponse{}, nil
}
