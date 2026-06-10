package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/domain/game"
)

func (g MultiGrpcServer) SubmitHPConsensus(ctx context.Context, request *protomulti.GameHPConsensusRequest) (*protomulti.GameHPConsensusResponse, error) {

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
