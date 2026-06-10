package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyCancelCommand struct {
	BasicCommand
}

func (c LobbyCancelCommand) CommandName() string {
	return "lobby.cancel"
}

func NewLobbyCancelCommand() *LobbyCancelCommand {
	return &LobbyCancelCommand{}
}

var (
	ErrFailedHandleCancelLobbyCommand = func(err error) error {
		return fmt.Errorf("failed handle cancel lobby command: %w", err)
	}
)

type LobbyCancelCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     port.SessionInstanceRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository port.LobbyParticipantRepository
}

func NewLobbyCancelCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, instances port.LobbyInstanceRepository, participants port.LobbyParticipantRepository) *LobbyCancelCommandHandler {
	return &LobbyCancelCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *LobbyCancelCommandHandler) Handle(ctx Context, cmd *LobbyCancelCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	instance, err := h.InstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := instance.CanCancel(ctx.PlayerID()); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	for _, participant := range participants {
		if !participant.HasPlayer() {
			continue
		}

		session, err := h.SessionRepository.QueryByID(ctx, participant.UserID)
		if err != nil {
			return ErrFailedHandleCancelLobbyCommand(err)
		}

		if session.LobbyID != instance.SysID {
			continue
		}

		session.LobbyID = uuid.Nil
		session.PartySlot = 0
		if err := h.SessionRepository.Update(ctx, session); err != nil {
			return ErrFailedHandleCancelLobbyCommand(err)
		}

		cmd.QueueEvent(lobby.NewParticipantDeletedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))
	}

	if err := h.ParticipantRepository.DeleteAllForLobby(ctx, instance.SysID); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := h.InstanceRepository.Delete(ctx, instance); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	cmd.QueueEvent(lobby.NewInstanceDeletedEvent(ctx, instance.SysID, instance.QuestID))
	return nil

}
