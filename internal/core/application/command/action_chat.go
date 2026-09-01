package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
)

type GameChatCommand struct {
	BasicCommand
	Message string
}

func (e *GameChatCommand) CommandName() string {
	return "action.chat"
}

func NewGameChatCommand(message string) *GameChatCommand {
	return &GameChatCommand{Message: message}
}

type GameChatCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewGameChatCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *GameChatCommandHandler {
	return &GameChatCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *GameChatCommandHandler) Handle(ctx Context, cmd *GameChatCommand) error {

	message, err := sanitiseChatMessage(cmd.Message)
	if err != nil {
		return err
	}

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerChatAction(current.GameID, current.PlayerID, message),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
