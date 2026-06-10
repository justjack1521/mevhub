package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) ParticipantLeave(ctx context.Context, request *protomulti.ParticipantLeaveRequest) (*protomulti.ParticipantLeaveResponse, error) {
	var cmd = command.NewParticipantLeaveCommand()

	if err := g.app.SubApplications.Lobby.Commands.ParticipantLeave.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.ParticipantLeaveResponse{}, nil
}
