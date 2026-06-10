package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
)

type SubmitHPConsensusCommand struct {
	BasicCommand
	Enemies []game.EnemyHP
}

func (c *SubmitHPConsensusCommand) CommandName() string {
	return "action.hp_consensus"
}

func NewSubmitHPConsensusCommand(enemies []game.EnemyHP) *SubmitHPConsensusCommand {
	return &SubmitHPConsensusCommand{Enemies: enemies}
}

type SubmitHPConsensusCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewSubmitHPConsensusCommandHandler(sessions port.SessionInstanceReadRepository, svr *server.GameServerHost) *SubmitHPConsensusCommandHandler {
	return &SubmitHPConsensusCommandHandler{SessionRepository: sessions, GameServerHost: svr}
}

func (h *SubmitHPConsensusCommandHandler) Handle(ctx Context, cmd *SubmitHPConsensusCommand) error {
	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if uuid.Equal(current.GameID, uuid.Nil) {
		return nil
	}

	var request = &server.GameActionRequest{
		GameID: current.GameID,
		Action: action.NewSubmitHPConsensusAction(current.GameID, current.PlayerID, cmd.Enemies),
	}

	h.GameServerHost.ActionChannel <- request

	return nil
}
