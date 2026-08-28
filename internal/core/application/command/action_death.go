package command

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
	"time"
)

// PlayerDeathCommand is the dying client's self-report: like the HP consensus,
// the server does not simulate combat, so it learns of a death from the client
// that suffered it. Identity comes from the session — a client can only report
// its own death.
type PlayerDeathCommand struct {
	BasicCommand
}

func (e *PlayerDeathCommand) CommandName() string {
	return "action.death"
}

func NewPlayerDeathCommand() *PlayerDeathCommand {
	return &PlayerDeathCommand{}
}

type PlayerDeathCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewPlayerDeathCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *PlayerDeathCommandHandler {
	return &PlayerDeathCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *PlayerDeathCommandHandler) Handle(ctx Context, cmd *PlayerDeathCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerDeathAction(current.GameID, current.PlayerID, time.Now().UTC()),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
