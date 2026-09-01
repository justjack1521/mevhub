package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) LobbyChat(ctx context.Context, request *protomulti.LobbyChatRequest) (*protomulti.LobbyChatResponse, error) {

	var cmd = command.NewLobbyChatCommand(request.Message)

	if err := g.app.SubApplications.Lobby.Commands.LobbyChat.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.LobbyChatResponse{}, nil

}
