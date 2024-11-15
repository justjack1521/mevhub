package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
)

type CancelLobbyCommand struct {
	BasicCommand
}

func (c CancelLobbyCommand) CommandName() string {
	return "lobby.cancel"
}

func NewCancelLobbyCommand() *CancelLobbyCommand {
	return &CancelLobbyCommand{}
}

var (
	ErrFailedHandleCancelLobbyCommand = func(err error) error {
		return fmt.Errorf("failed handle cancel lobby command: %w", err)
	}
)

type CancelLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    lobby.InstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewCancelLobbyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances lobby.InstanceRepository, participants lobby.ParticipantRepository) *CancelLobbyCommandHandler {
	return &CancelLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *CancelLobbyCommandHandler) Handle(ctx Context, cmd *CancelLobbyCommand) error {

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

	if err := h.InstanceRepository.Delete(ctx, current.LobbyID); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	cmd.QueueEvent(lobby.NewInstanceDeletedEvent(ctx, current.LobbyID))

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx, current.LobbyID)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if err := h.ParticipantRepository.Delete(ctx, participant); err != nil {
			return err
		}
		cmd.QueueEvent(lobby.NewParticipantDeletedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))
	}

	return nil

}
