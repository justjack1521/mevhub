package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
)

// GameCatchUp resends the notifications a client missed. The backlog itself
// leaves over the client's normal notification channel, in order, so only the
// range it covers comes back on this call.
func (g MultiGrpcServer) GameCatchUp(ctx context.Context, request *protomulti.GameCatchUpRequest) (*protomulti.GameCatchUpResponse, error) {

	cmd := command.NewGameCatchUpCommand(request.GetLastSequence())

	if err := g.app.SubApplications.Game.Commands.GameCatchUp.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameCatchUpResponse{
		FromSequence:    cmd.Result.FromSequence,
		ToSequence:      cmd.Result.ToSequence,
		TurnRemainingMs: cmd.Result.TurnRemainingMs,
	}, nil

}
