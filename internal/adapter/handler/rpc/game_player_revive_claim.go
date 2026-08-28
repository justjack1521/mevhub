package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/command"
)

func (g MultiGrpcServer) PlayerReviveClaim(ctx context.Context, request *protomulti.GamePlayerReviveClaimRequest) (*protomulti.GamePlayerReviveClaimResponse, error) {

	// An empty or malformed target parses to uuid.Nil, which the command
	// resolves to the caller — claiming a self-revive.
	var cmd = command.NewPlayerReviveClaimCommand(uuid.FromStringOrNil(request.TargetPlayerId))

	if err := g.app.SubApplications.Game.Commands.PlayerReviveClaim.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	// A denial is a successful RPC carrying "no": the deny-reason ordinals are
	// kept in sync with the domain enum, so this is a straight cast.
	return &protomulti.GamePlayerReviveClaimResponse{
		Granted:     cmd.Result.Granted,
		DenyReason:  protomulti.ReviveClaimDenyReason(cmd.Result.Reason),
		RemainingMs: cmd.Result.Remaining.Milliseconds(),
	}, nil

}
