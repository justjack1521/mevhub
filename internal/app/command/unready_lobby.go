package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type UnreadyLobbyCommand struct {
	BasicCommand
	LobbyID uuid.UUID
}

func (c UnreadyLobbyCommand) CommandName() string {
	return "lobby.unready"
}

func NewUnreadyLobbyCommand(id uuid.UUID) *UnreadyLobbyCommand {
	return &UnreadyLobbyCommand{LobbyID: id}
}

type UnreadyLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    lobby.InstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewUnreadyLobbyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances lobby.InstanceRepository, participants lobby.ParticipantRepository) *UnreadyLobbyCommandHandler {
	return &UnreadyLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *UnreadyLobbyCommandHandler) Handle(ctx Context, cmd *UnreadyLobbyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, current.LobbyID, current.PartySlot)
	if err != nil {
		return err
	}

	if err := participant.SetReady(ctx.PlayerID(), false); err != nil {
		return err
	}

	if err := h.ParticipantRepository.Update(ctx, participant); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantUnreadyEvent(ctx, current.UserID, current.LobbyID, current.PartySlot))

	return nil

}
