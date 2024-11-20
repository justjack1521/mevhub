package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
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
	SessionRepository     port.SessionInstanceReadRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository port.LobbyParticipantRepository
}

func NewLobbyCancelCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceReadRepository, instances port.LobbyInstanceRepository, participants port.LobbyParticipantRepository) *LobbyCancelCommandHandler {
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

	if err := h.InstanceRepository.Delete(ctx, instance); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	cmd.QueueEvent(lobby.NewInstanceDeletedEvent(ctx, instance.SysID, instance.QuestID))
	return nil

}
