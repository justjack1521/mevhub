package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type ParticipantUnreadyCommand struct {
	BasicCommand
	LobbyID uuid.UUID
}

func (c ParticipantUnreadyCommand) CommandName() string {
	return "participant.unready"
}

func NewParticipantUnreadyCommand(id uuid.UUID) *ParticipantUnreadyCommand {
	return &ParticipantUnreadyCommand{LobbyID: id}
}

type ParticipantUnreadyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository port.LobbyParticipantRepository
}

func NewParticipantUnreadyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances port.LobbyInstanceRepository, participants port.LobbyParticipantRepository) *ParticipantUnreadyCommandHandler {
	return &ParticipantUnreadyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *ParticipantUnreadyCommandHandler) Handle(ctx Context, cmd *ParticipantUnreadyCommand) error {

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
