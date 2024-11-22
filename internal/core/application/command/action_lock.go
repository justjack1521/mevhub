package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
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
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerLockAction(current.GameID, current.LobbyID, current.PlayerID),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
