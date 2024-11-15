package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) CreateLobby(context context.Context, request *protomulti.CreateLobbyRequest) (*protomulti.CreateLobbyResponse, error) {
	quest, err := uuid.FromString(request.QuestId)
	if err != nil {
		return nil, err
	}

	var options = command.CreateLobbyOptions{
		MinimumPlayerLevel: int(request.MinPlayerLevel),
		Restrictions:       nil,
	}

	var cmd = command.NewCreateLobbyCommand(quest, int(request.DeckIndex), request.Comment, options)

	if err := g.app.SubApplications.Lobby.Commands.CreateLobby.Handle(g.NewCommandContext(context), cmd); err != nil {
		return nil, err
	}

	return &protomulti.CreateLobbyResponse{
		LobbyId: cmd.LobbyID.String(),
		PartyId: cmd.PartyID,
	}, nil
}
