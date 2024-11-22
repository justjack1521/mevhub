package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
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
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewReadyPlayerCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *ReadyPlayerCommandHandler {
	return &ReadyPlayerCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *ReadyPlayerCommandHandler) Handle(ctx Context, cmd *ReadyPlayerCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  game.NewPlayerReadyAction(current.GameID, current.LobbyID, current.PlayerID),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
