package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
)

// PlayerReviveCommand revives a dead player. The target may be another player
// or the caller themselves; a nil target means self-revive. The reviver is
// always the session's player.
type PlayerReviveCommand struct {
	BasicCommand
	TargetPlayerID uuid.UUID
}

func (e *PlayerReviveCommand) CommandName() string {
	return "action.revive"
}

func NewPlayerReviveCommand(target uuid.UUID) *PlayerReviveCommand {
	return &PlayerReviveCommand{TargetPlayerID: target}
}

type PlayerReviveCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewPlayerReviveCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *PlayerReviveCommandHandler {
	return &PlayerReviveCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *PlayerReviveCommandHandler) Handle(ctx Context, cmd *PlayerReviveCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var target = cmd.TargetPlayerID
	if uuid.Equal(target, uuid.Nil) {
		target = current.PlayerID
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerReviveAction(current.GameID, target, current.PlayerID),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
