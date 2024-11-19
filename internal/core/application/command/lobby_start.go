package command

import (
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyStartCommand struct {
	BasicCommand
}

func (c LobbyStartCommand) CommandName() string {
	return "lobby.start"
}

func NewLobbyStartCommand() *LobbyStartCommand {
	return &LobbyStartCommand{}
}

type LobbyStartCommandHandler struct {
	SessionRepository          port.SessionInstanceReadRepository
	LobbyInstanceRepository    port.LobbyInstanceRepository
	LobbyParticipantRepository port.LobbyParticipantReadRepository
}

func NewLobbyStartCommandHandler(sessions port.SessionInstanceReadRepository, lobbies port.LobbyInstanceRepository, participants port.LobbyParticipantReadRepository) *LobbyStartCommandHandler {
	return &LobbyStartCommandHandler{
		SessionRepository:          sessions,
		LobbyInstanceRepository:    lobbies,
		LobbyParticipantRepository: participants,
	}
}

func (h *LobbyStartCommandHandler) Handle(ctx Context, cmd *LobbyStartCommand) error {

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

	return nil

}
