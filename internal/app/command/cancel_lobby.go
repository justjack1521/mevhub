package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type CancelLobbyCommand struct {
}

func (c CancelLobbyCommand) CommandName() string {
	return "cancel.lobby"
}

func NewCancelLobbyCommand() CancelLobbyCommand {
	return CancelLobbyCommand{}
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

func (h *CancelLobbyCommandHandler) Handle(ctx *Context, cmd CancelLobbyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx.Context, ctx.ClientID)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	instance, err := h.InstanceRepository.QueryByID(ctx.Context, current.LobbyID)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := instance.CanCancel(ctx.ClientID); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := h.InstanceRepository.Delete(ctx.Context, current.LobbyID); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	h.EventPublisher.Notify(lobby.NewInstanceDeletedEvent(ctx.Context, current.LobbyID))

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx.Context, current.LobbyID)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if err := h.ParticipantRepository.Delete(ctx.Context, participant); err != nil {
			return err
		}
		h.EventPublisher.Notify(lobby.NewParticipantDeletedEvent(ctx.Context, participant.ClientID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))
	}

	return nil

}
