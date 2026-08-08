package command

import (
	"errors"

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
		// The slot is gone entirely — the lobby was cancelled or started. Clear
		// the session so the player is not stranded pointing at a dead lobby.
		if errors.Is(err, port.ErrParticipantNotFound) {
			current.LobbyID = uuid.Nil
			current.PartySlot = 0
			return h.SessionRepository.Update(ctx, current)
		}
		return err
	}

	var (
		user    = participant.UserID
		player  = participant.PlayerID
		lobbyID = participant.LobbyID
		slot    = participant.PlayerSlot
	)

	// A stale session can point at a slot someone else has since taken; clear the
	// session but leave the occupant alone.
	if !uuid.Equal(player, ctx.PlayerID()) {
		current.LobbyID = uuid.Nil
		current.PartySlot = 0
		return h.SessionRepository.Update(ctx, current)
	}

	// Slots are pre-created empty at lobby creation and joining mutates them in
	// place, so vacate the slot rather than deleting it — deleting removes the
	// placeholder the next join depends on.
	participant.RemovePlayer()
	if err := h.ParticipantRepository.Create(ctx, participant); err != nil {
		return err
	}

	current.LobbyID = uuid.Nil
	current.PartySlot = 0
	if err := h.SessionRepository.Update(ctx, current); err != nil {
		return err
	}

	if err := h.ListenerRepository.DeleteListener(ctx, lobbyID, ctx.UserID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantDeletedEvent(ctx, user, player, lobbyID, slot))

	return nil

}
