package rpc

import (
	"context"
	"errors"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/domain/game"
)

func (g MultiGrpcServer) SubmitHPConsensus(ctx context.Context, request *protomulti.GameHPConsensusRequest) (*protomulti.GameHPConsensusResponse, error) {

	// A report with no enemies carries no information and, unvalidated, an
	// all-empty round would read as a wipe and end the game for everyone.
	if len(request.Enemies) == 0 {
		return nil, errors.New("hp consensus report contains no enemies")
	}

	enemies := make([]game.EnemyHP, len(request.Enemies))
	for i, e := range request.Enemies {
		enemies[i] = game.EnemyHP{
			EnemyIndex: int(e.EnemyIndex),
			HP:         int(e.Hp),
		}
	}

	cmd := command.NewSubmitHPConsensusCommand(enemies)

	if err := g.app.SubApplications.Game.Commands.SubmitHPConsensus.Handle(g.NewCommandContext(ctx), cmd); err != nil {
		return nil, err
	}

	return &protomulti.GameHPConsensusResponse{}, nil

}
