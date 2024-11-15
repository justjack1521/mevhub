package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
)

type DequeueActionCommand struct {
	BasicCommand
}

func (e *DequeueActionCommand) CommandName() string {
	return "action.dequeue"
}

func NewDequeueActionCommand() *DequeueActionCommand {
	return &DequeueActionCommand{}
}

type DequeueActionCommandHandler struct {
	SessionRepository session.InstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewDequeueActionCommandHandler(sessions session.InstanceReadRepository, server *server.GameServerHost) *DequeueActionCommandHandler {
	return &DequeueActionCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *DequeueActionCommandHandler) Handle(ctx Context, cmd *DequeueActionCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		InstanceID: current.LobbyID,
		Action: &game.PlayerDequeueAction{
			InstanceID: current.LobbyID,
			PlayerID:   current.PlayerID,
		},
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
