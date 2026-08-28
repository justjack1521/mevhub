package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) PlayerRevive(ctx context.Context, request *protomulti.GamePlayerReviveRequest) (*protomulti.GamePlayerReviveResponse, error) {

	// An empty or malformed target parses to uuid.Nil, which the command
	// resolves to the caller — a self-revive.
	var cmd = command.NewPlayerReviveCommand(uuid.FromStringOrNil(request.TargetPlayerId))

	if err := g.app.SubApplications.Game.Commands.PlayerRevive.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GamePlayerReviveResponse{}, nil

}
