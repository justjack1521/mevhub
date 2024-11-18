package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantFindResponse(ctx context.Context, request *protomulti.ParticipantFindRequest) (*protomulti.ParticipantFindResponse, error) {

	id, err := uuid.FromString(request.QuestId)
	if err != nil {
		return nil, err
	}

	var cmd = command.NewParticipantFindCommand(id, int(request.DeckIndex), request.UseStamina)

	if err := g.app.SubApplications.Lobby.Commands.ParticipantFind.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantFindResponse{}, nil

}
