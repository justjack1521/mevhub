package command

import (
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type StartLobbyCommand struct {
	BasicCommand
}

func (c StartLobbyCommand) CommandName() string {
	return "lobby.start"
}

func NewStartLobbyCommand() *StartLobbyCommand {
	return &StartLobbyCommand{}
}

type StartLobbyCommandHandler struct {
	SessionRepository  session.InstanceReadRepository
	InstanceRepository lobby.InstanceRepository
}

func NewStartLobbyCommandHandler(sessions session.InstanceReadRepository, lobbies lobby.InstanceRepository) *StartLobbyCommandHandler {
	return &StartLobbyCommandHandler{SessionRepository: sessions, InstanceRepository: lobbies}
}

func (h *StartLobbyCommandHandler) Handle(ctx Context, cmd *StartLobbyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	instance, err := h.InstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		return err
	}

	if err := instance.StartLobby(ctx.PlayerID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewInstanceStartedEvent(ctx, instance.SysID))

	return nil

}
