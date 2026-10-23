package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
	"time"
)

type GameStampCommand struct {
	BasicCommand
	StampID uuid.UUID
}

func (e *GameStampCommand) CommandName() string {
	return "action.stamp"
}

func NewGameStampCommand(id uuid.UUID) *GameStampCommand {
	return &GameStampCommand{StampID: id}
}

type GameStampCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewGameStampCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *GameStampCommandHandler {
	return &GameStampCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *GameStampCommandHandler) Handle(ctx Context, cmd *GameStampCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerStampAction(current.GameID, current.PlayerID, cmd.StampID, time.Now().UTC()),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
