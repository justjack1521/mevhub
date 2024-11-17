package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) JoinLobby(ctx context.Context, request *protomulti.JoinLobbyRequest) (*protomulti.JoinLobbyResponse, error) {
	id, err := uuid.FromString(request.LobbyId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewParticipantJoinCommand(id, int(request.DeckIndex), int(request.SlotIndex), request.UseStamina, request.FromInvite)

	if err := g.app.SubApplications.Lobby.Commands.ParticipantJoin.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.JoinLobbyResponse{}, nil
}
