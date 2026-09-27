package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/domain/lobby"
)

func (g MultiGrpcServer) LobbyCreate(context context.Context, request *protomulti.LobbyCreateRequest) (*protomulti.LobbyCreateResponse, error) {
	quest, err := uuid.FromString(request.QuestId)
	if err != nil {
		return nil, err
	}

	var restrictions = make([]lobby.PlayerSlotRestriction, len(request.Restrictions))
	for index, value := range request.Restrictions {
		restrictions[index] = lobby.PlayerSlotRestriction{
			Index:           int(value.SlotIndex),
			RoleRestriction: uuid.FromStringOrNil(value.RoleRestriction),
			InviteOnly:      value.InviteOnly,
			Locked:          value.Locked,
			BotControl:      value.Bot,
		}
	}

	var options = command.CreateLobbyOptions{
		MinimumPlayerLevel: int(request.MinPlayerLevel),
		Restrictions:       restrictions,
	}

	var cmd = command.NewLobbyCreateCommand(quest, int(request.DeckIndex), request.Comment, options)

	if err := g.app.SubApplications.Lobby.Commands.LobbyCreate.Handle(g.NewCommandContext(context), cmd); err != nil {
		return nil, err
	}

	return &protomulti.LobbyCreateResponse{
		LobbyId: cmd.LobbyID.String(),
		PartyId: cmd.PartyID,
	}, nil
}
