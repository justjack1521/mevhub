package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantUnwatch(ctx context.Context, request *protomulti.ParticipantUnwatchRequest) (*protomulti.ParticipantUnwatchResponse, error) {
	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewUnwatchLobbyCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.ParticipantUnwatch.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantUnwatchResponse{}, nil
}
