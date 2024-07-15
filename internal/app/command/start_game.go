package command

import (
	uuid "github.com/satori/go.uuid"
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
	SessionRepository          session.InstanceReadRepository
	LobbyInstanceRepository    lobby.InstanceRepository
	LobbyParticipantRepository lobby.ParticipantReadRepository
	GameInstanceFactory        *game.InstanceFactory
	GameInstanceRepository     game.InstanceWriteRepository
	GameParticipantFactory     *game.PlayerParticipantFactory
	GameParticipantRepository  game.PlayerParticipantWriteRepository
}

func NewStartLobbyCommandHandler(sessions session.InstanceReadRepository, lobbies lobby.InstanceRepository, participants lobby.ParticipantReadRepository, factory *game.InstanceFactory, games game.InstanceWriteRepository, gameParticipantFactory *game.PlayerParticipantFactory, gameParticipantRepository game.PlayerParticipantWriteRepository) *StartLobbyCommandHandler {
	return &StartLobbyCommandHandler{SessionRepository: sessions, LobbyInstanceRepository: lobbies, LobbyParticipantRepository: participants, GameInstanceFactory: factory, GameInstanceRepository: games, GameParticipantFactory: gameParticipantFactory, GameParticipantRepository: gameParticipantRepository}
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

	participants, err := h.LobbyParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if uuid.Equal(participant.PlayerID, uuid.Nil) {
			continue
		}
		part, err := h.GameParticipantFactory.Create(ctx, participant)
		if err != nil {
			return err
		}
		if err := h.GameParticipantRepository.Create(ctx, inst.SysID, participant.PlayerSlot, part); err != nil {
			return err
		}

	}

	return nil

}
