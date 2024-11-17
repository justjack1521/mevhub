package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ReadyLobby(ctx context.Context, request *protomulti.ReadyLobbyRequest) (*protomulti.ReadyLobbyResponse, error) {

	lobby, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewReadyParticipantCommand(lobby, int(request.DeckIndex))

	if err := g.app.SubApplications.Lobby.Commands.ParticipantReady.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ReadyLobbyResponse{}, nil
}
