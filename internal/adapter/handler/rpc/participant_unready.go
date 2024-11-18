package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantUnready(ctx context.Context, request *protomulti.ParticipantUnreadyRequest) (*protomulti.ParticipantUnreadyResponse, error) {
	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewParticipantUnreadyCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.ParticipantUnready.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantUnreadyResponse{}, nil
}
