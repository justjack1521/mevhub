package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type ParticipantLeaveCommand struct {
	BasicCommand
}

func (c ParticipantLeaveCommand) CommandName() string {
	return "participant.leave"
}

func NewParticipantLeaveCommand() *ParticipantLeaveCommand {
	return &ParticipantLeaveCommand{}
}

type ParticipantLeaveCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     port.SessionInstanceRepository
	ParticipantRepository port.LobbyParticipantRepository
	ListenerRepository    lobby.NotificationListenerWriteRepository
}

func NewParticipantLeaveCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, participants port.LobbyParticipantRepository, listeners lobby.NotificationListenerWriteRepository) *ParticipantLeaveCommandHandler {
	return &ParticipantLeaveCommandHandler{EventPublisher: publisher, SessionRepository: sessions, ParticipantRepository: participants, ListenerRepository: listeners}
}

func (h *ParticipantLeaveCommandHandler) Handle(ctx Context, cmd *ParticipantLeaveCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, current.LobbyID, current.PartySlot)
	if err != nil {
		return err
	}

	if err := h.ParticipantRepository.Delete(ctx, participant); err != nil {
		return err
	}

	current.LobbyID = uuid.Nil
	current.PartySlot = 0
	if err := h.SessionRepository.Update(ctx, current); err != nil {
		return err
	}

	if err := h.ListenerRepository.DeleteListener(ctx, participant.LobbyID, ctx.UserID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantDeletedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))

	return nil

}
