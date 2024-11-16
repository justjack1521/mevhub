package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type UnreadyParticipantCommand struct {
	BasicCommand
	LobbyID uuid.UUID
}

func (c UnreadyParticipantCommand) CommandName() string {
	return "participant.unready"
}

func NewUnreadyParticipantCommand(id uuid.UUID) *UnreadyParticipantCommand {
	return &UnreadyParticipantCommand{LobbyID: id}
}

type UnreadyParticipantCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewUnreadyParticipantCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances port.LobbyInstanceRepository, participants lobby.ParticipantRepository) *UnreadyParticipantCommandHandler {
	return &UnreadyParticipantCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *UnreadyParticipantCommandHandler) Handle(ctx Context, cmd *UnreadyParticipantCommand) error {

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
