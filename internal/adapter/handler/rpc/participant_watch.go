package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantWatch(ctx context.Context, request *protomulti.ParticipantWatchRequest) (*protomulti.ParticipantWatchResponse, error) {
	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewWatchLobbyCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.ParticipantWatch.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantWatchResponse{}, nil
}
