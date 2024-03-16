package command

import (
	"errors"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type LeaveLobbyCommand struct {
}

func (c LeaveLobbyCommand) CommandName() string {
	return "leave.lobby"
}

func NewLeaveLobbyCommand() LeaveLobbyCommand {
	return LeaveLobbyCommand{}
}

type LeaveLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewLeaveLobbyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceRepository, participants lobby.ParticipantRepository) *LeaveLobbyCommandHandler {
	return &LeaveLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, ParticipantRepository: participants}
}

func (h *LeaveLobbyCommandHandler) Handle(ctx *Context, cmd LeaveLobbyCommand) error {

	if ctx.Session.CanLeaveLobby() == false {
		return errors.New("player has already left lobby")
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx.Context, ctx.Session.LobbyID, ctx.Session.PartySlot)
	if err != nil {
		return err
	}

	if err := h.ParticipantRepository.Delete(ctx.Context, participant); err != nil {
		return err
	}

	h.EventPublisher.Notify(lobby.NewParticipantDeletedEvent(ctx.Context, participant.ClientID, participant.PlayerID, participant.LobbyID, participant.PlayerSlot))

	return nil

}
