package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
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
}

func NewParticipantLeaveCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, participants port.LobbyParticipantRepository) *ParticipantLeaveCommandHandler {
	return &ParticipantLeaveCommandHandler{EventPublisher: publisher, SessionRepository: sessions, ParticipantRepository: participants}
}

func (h *ParticipantLeaveCommandHandler) Handle(ctx Context, cmd *ParticipantLeaveCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, current.LobbyID, current.PartySlot)
	if err != nil {
		return err
	}

	if err := h.ParticipantRepository.Delete(ctx, participant); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantDeletedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))

	return nil

}
