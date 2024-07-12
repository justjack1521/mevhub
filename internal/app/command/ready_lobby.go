package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type ReadyLobbyCommand struct {
	LobbyID   uuid.UUID
	DeckIndex int
}

func (c ReadyLobbyCommand) CommandName() string {
	return "ready.lobby"
}

func NewReadyLobbyCommand(id uuid.UUID, deck int) ReadyLobbyCommand {
	return ReadyLobbyCommand{LobbyID: id, DeckIndex: deck}
}

type ReadyLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     session.InstanceReadRepository
	InstanceRepository    lobby.InstanceRepository
	ParticipantRepository lobby.ParticipantRepository
}

func NewReadyLobbyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, instances lobby.InstanceRepository, participants lobby.ParticipantRepository) *ReadyLobbyCommandHandler {
	return &ReadyLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *ReadyLobbyCommandHandler) Handle(ctx Context, cmd ReadyLobbyCommand) error {

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

	h.EventPublisher.Notify(lobby.NewParticipantReadyEvent(ctx, ctx.UserID(), current.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	if participant.DeckIndex != cmd.DeckIndex {
		h.EventPublisher.Notify(lobby.NewParticipantDeckChangeEvent(ctx, ctx.UserID(), ctx.PlayerID(), current.LobbyID, cmd.DeckIndex, participant.PlayerSlot))
	}

	return nil

}
