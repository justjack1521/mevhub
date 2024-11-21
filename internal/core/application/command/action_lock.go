package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type LockActionCommand struct {
	BasicCommand
}

func (e *LockActionCommand) CommandName() string {
	return "action.lock"
}

func NewLockActionCommand() *LockActionCommand {
	return &LockActionCommand{}
}

type LockActionCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewLockActionCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *LockActionCommandHandler {
	return &LockActionCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *LockActionCommandHandler) Handle(ctx Context, cmd *LockActionCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		PartyID: current.LobbyID,
		Action: &game.PlayerLockAction{
			InstanceID: current.LobbyID,
			PlayerID:   current.PlayerID,
		},
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
