package ports

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/command"
)

func (g MultiGrpcServer) ReadyLobby(ctx context.Context, request *protomulti.ReadyLobbyRequest) (*protomulti.ReadyLobbyResponse, error) {

	lobby, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewReadyLobbyCommand(lobby, int(request.DeckIndex))

	if err := g.app.SubApplications.Lobby.Commands.ReadyLobby.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ReadyLobbyResponse{}, nil
}
