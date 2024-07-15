package command

import (
	"mevhub/internal/domain/game"
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
	SessionRepository       session.InstanceReadRepository
	LobbyInstanceRepository lobby.InstanceRepository
	GameInstanceFactory     *game.InstanceFactory
	GameInstanceRepository  game.InstanceRepository
}

func NewStartLobbyCommandHandler(session session.InstanceReadRepository, lobbies lobby.InstanceRepository, factory *game.InstanceFactory, games game.InstanceRepository) *StartLobbyCommandHandler {
	return &StartLobbyCommandHandler{SessionRepository: session, LobbyInstanceRepository: lobbies, GameInstanceFactory: factory, GameInstanceRepository: games}
}

func (h *StartLobbyCommandHandler) Handle(ctx Context, cmd *StartLobbyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	instance, err := h.LobbyInstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		return err
	}

	if err := instance.StartLobby(ctx.PlayerID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewInstanceStartedEvent(ctx, instance.SysID))

	if err := h.LobbyInstanceRepository.Create(ctx, instance); err != nil {
		return err
	}

	inst, err := h.GameInstanceFactory.Create(instance)
	if err != nil {
		return err
	}

	if err := h.GameInstanceRepository.Create(ctx, inst); err != nil {
		return err
	}

	return nil

}
