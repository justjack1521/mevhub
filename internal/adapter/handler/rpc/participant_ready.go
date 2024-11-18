package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantReady(ctx context.Context, request *protomulti.ParticipantReadyRequest) (*protomulti.ParticipantReadyResponse, error) {

	lobby, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewParticipantReadyCommand(lobby, int(request.DeckIndex))

	if err := g.app.SubApplications.Lobby.Commands.ParticipantReady.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantReadyResponse{}, nil
}
