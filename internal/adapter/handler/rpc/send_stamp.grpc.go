package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) SendStamp(ctx context.Context, request *protomulti.SendStampRequest) (*protomulti.SendStampResponse, error) {
	id, err := uuid.FromString(request.StampId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewLobbyStampCommand(id)

	if err := g.app.SubApplications.Lobby.Commands.LobbyStamp.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.SendStampResponse{}, nil
}
