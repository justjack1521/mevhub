package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
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
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewDequeueActionCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *DequeueActionCommandHandler {
	return &DequeueActionCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *DequeueActionCommandHandler) Handle(ctx Context, cmd *DequeueActionCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  game.NewPlayerDequeueAction(current.GameID, current.LobbyID, current.PlayerID),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
