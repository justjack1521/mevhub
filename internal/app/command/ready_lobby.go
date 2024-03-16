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

func (h *ReadyLobbyCommandHandler) Handle(ctx *Context, cmd ReadyLobbyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx.Context, ctx.ClientID)
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx.Context, current.LobbyID, current.PartySlot)
	if err != nil {
		return err
	}

	if err := participant.SetReady(ctx.Session.PlayerID, true); err != nil {
		return err
	}

	if err := participant.SetDeckIndex(ctx.Session.PlayerID, cmd.DeckIndex); err != nil {
		return err
	}

	if err := h.ParticipantRepository.Update(ctx.Context, participant); err != nil {
		return err
	}

	h.EventPublisher.Notify(lobby.NewParticipantReadyEvent(ctx.Context, ctx.ClientID, current.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	if participant.DeckIndex != cmd.DeckIndex {
		h.EventPublisher.Notify(lobby.NewParticipantDeckChangeEvent(ctx.Context, ctx.ClientID, ctx.Session.PlayerID, current.LobbyID, cmd.DeckIndex, participant.PlayerSlot))
	}

	return nil

}
