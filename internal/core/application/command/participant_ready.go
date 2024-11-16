package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type ReadyParticipantCommand struct {
	BasicCommand
	LobbyID   uuid.UUID
	DeckIndex int
}

func (c ReadyParticipantCommand) CommandName() string {
	return "participant.ready"
}

func NewReadyParticipantCommand(id uuid.UUID, deck int) *ReadyParticipantCommand {
	return &ReadyParticipantCommand{LobbyID: id, DeckIndex: deck}
}

type ReadyLobbyParticipantHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewReadyParticipantCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances port.LobbyInstanceRepository, participants lobby.ParticipantRepository) *ReadyLobbyParticipantHandler {
	return &ReadyLobbyParticipantHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *ReadyLobbyParticipantHandler) Handle(ctx Context, cmd *ReadyParticipantCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, current.LobbyID, current.PartySlot)
	if err != nil {
		return err
	}

	if err := participant.SetReady(ctx.PlayerID(), true); err != nil {
		return err
	}

	if err := participant.SetDeckIndex(ctx.PlayerID(), cmd.DeckIndex); err != nil {
		return err
	}

	if err := h.ParticipantRepository.Update(ctx, participant); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantReadyEvent(ctx, ctx.UserID(), current.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	if participant.DeckIndex != cmd.DeckIndex {
		cmd.QueueEvent(lobby.NewParticipantDeckChangeEvent(ctx, ctx.UserID(), ctx.PlayerID(), current.LobbyID, cmd.DeckIndex, participant.PlayerSlot))
	}

	return nil

}
