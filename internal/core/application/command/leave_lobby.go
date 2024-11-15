package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
)

type LeaveLobbyCommand struct {
	BasicCommand
}

func (c LeaveLobbyCommand) CommandName() string {
	return "lobby.leave"
}

func NewLeaveLobbyCommand() *LeaveLobbyCommand {
	return &LeaveLobbyCommand{}
}

type LeaveLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewLeaveLobbyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceRepository, participants lobby.ParticipantRepository) *LeaveLobbyCommandHandler {
	return &LeaveLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, ParticipantRepository: participants}
}

func (h *LeaveLobbyCommandHandler) Handle(ctx Context, cmd *LeaveLobbyCommand) error {

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
