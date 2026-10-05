package command

import (
	"errors"
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyCancelCommand struct {
	BasicCommand
}

func (c LobbyCancelCommand) CommandName() string {
	return "lobby.cancel"
}

func NewLobbyCancelCommand() *LobbyCancelCommand {
	return &LobbyCancelCommand{}
}

var (
	ErrFailedHandleCancelLobbyCommand = func(err error) error {
		return fmt.Errorf("failed handle cancel lobby command: %w", err)
	}
)

type LobbyCancelCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     port.SessionInstanceRepository
	InstanceRepository    port.LobbyInstanceRepository
	ParticipantRepository port.LobbyParticipantRepository
}

func NewLobbyCancelCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, instances port.LobbyInstanceRepository, participants port.LobbyParticipantRepository) *LobbyCancelCommandHandler {
	return &LobbyCancelCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants}
}

func (h *LobbyCancelCommandHandler) Handle(ctx Context, cmd *LobbyCancelCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	instance, err := h.InstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		if errors.Is(err, lobby.ErrLobbyInstanceNotFound) && !uuid.Equal(current.LobbyID, uuid.Nil) {
			return h.cancelExpired(ctx, cmd, current.LobbyID)
		}
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := instance.CanCancel(ctx.PlayerID()); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	for _, participant := range participants {
		if !participant.HasPlayer() {
			continue
		}

		session, err := h.SessionRepository.QueryByID(ctx, participant.UserID)
		if err != nil {
			return ErrFailedHandleCancelLobbyCommand(err)
		}

		if session.LobbyID != instance.SysID {
			continue
		}

		session.LobbyID = uuid.Nil
		session.PartySlot = 0
		if err := h.SessionRepository.Update(ctx, session); err != nil {
			return ErrFailedHandleCancelLobbyCommand(err)
		}
	}

	// Deliberately no per-participant deleted event here. Its handlers drop the
	// player from the listener set, and events publish in queue order — draining
	// the listeners first would leave InstanceDeletedEvent with nobody to send
	// the cancel notification to. The instance event covers the whole teardown.

	if err := h.ParticipantRepository.DeleteAllForLobby(ctx, instance.SysID); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	if err := h.InstanceRepository.Delete(ctx, instance); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	cmd.QueueEvent(lobby.NewInstanceDeletedEvent(ctx, instance.SysID, instance.QuestID))
	return nil

}

// cancelExpired finishes the teardown of a lobby whose instance already expired out of
// Redis. Expiry emits nothing, so this is the first anyone hears of it: free the
// participants' sessions, and still publish the deletion so the remaining listeners are
// sent CANCEL and the lobby's notification channel is retired. Without the instance there
// is no host to check against, but there is also nothing left to protect.
func (h *LobbyCancelCommandHandler) cancelExpired(ctx Context, cmd *LobbyCancelCommand, id uuid.UUID) error {

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx, id)
	if err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	var users = make([]uuid.UUID, 0, len(participants)+1)
	users = append(users, ctx.UserID())
	for _, participant := range participants {
		if participant.HasPlayer() {
			users = append(users, participant.UserID)
		}
	}

	for _, user := range users {
		session, err := h.SessionRepository.QueryByID(ctx, user)
		if err != nil {
			continue
		}
		if session.LobbyID != id {
			continue
		}
		session.LobbyID = uuid.Nil
		session.PartySlot = 0
		if err := h.SessionRepository.Update(ctx, session); err != nil {
			return ErrFailedHandleCancelLobbyCommand(err)
		}
	}

	if err := h.ParticipantRepository.DeleteAllForLobby(ctx, id); err != nil {
		return ErrFailedHandleCancelLobbyCommand(err)
	}

	// The quest went with the instance. Subscribers keyed on it (search and queue
	// clean-up) cannot act without it; the notification channeler only needs the id.
	cmd.QueueEvent(lobby.NewInstanceDeletedEvent(ctx, id, uuid.Nil))
	return nil

}
