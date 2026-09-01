package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) GameChat(ctx context.Context, request *protomulti.GameChatRequest) (*protomulti.GameChatResponse, error) {

	var cmd = command.NewGameChatCommand(request.Message)

	if err := g.app.SubApplications.Game.Commands.GameChat.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameChatResponse{}, nil

}
