package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
)

type ReadyPlayerCommand struct {
	BasicCommand
}

func (e *ReadyPlayerCommand) CommandName() string {
	return "action.ready"
}

func NewReadyPlayerCommand() *ReadyPlayerCommand {
	return &ReadyPlayerCommand{}
}

type ReadyPlayerCommandHandler struct {
	SessionRepository session.InstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewReadyPlayerCommandHandler(sessions session.InstanceReadRepository, server *server.GameServerHost) *ReadyPlayerCommandHandler {
	return &ReadyPlayerCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *ReadyPlayerCommandHandler) Handle(ctx Context, cmd *ReadyPlayerCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		InstanceID: current.LobbyID,
		Action: &game.PlayerReadyAction{
			InstanceID: current.LobbyID,
			PlayerID:   current.PlayerID,
		},
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
