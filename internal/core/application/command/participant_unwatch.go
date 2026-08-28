package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type UnwatchLobbyCommand struct {
	BasicCommand
	LobbyID uuid.UUID
}

func (c UnwatchLobbyCommand) CommandName() string {
	return "lobby.unwatch"
}

func NewUnwatchLobbyCommand(lobby uuid.UUID) *UnwatchLobbyCommand {
	return &UnwatchLobbyCommand{LobbyID: lobby}
}

type UnwatchLobbyCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceReadRepository
}

func NewUnwatchLobbyCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceReadRepository) *UnwatchLobbyCommandHandler {
	return &UnwatchLobbyCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *UnwatchLobbyCommandHandler) Handle(ctx Context, cmd *UnwatchLobbyCommand) error {

	// Participants and watchers share one listener set. A participant whose
	// session points at this lobby must not be able to unsubscribe themselves
	// from their own lobby's notifications by calling unwatch, so the session
	// arbitrates: skip on match, never error. A missing session cannot mark a
	// participant, so the removal proceeds — an expired watcher may still
	// clean up their listener.
	exists, err := h.SessionRepository.Exists(ctx, ctx.UserID())
	if err != nil {
		return err
	}
	if exists {
		current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
		if err != nil {
			return err
		}
		if uuid.Equal(current.LobbyID, cmd.LobbyID) {
			return nil
		}
	}

	cmd.QueueEvent(lobby.NewWatcherRemovedEvent(ctx, cmd.LobbyID, ctx.UserID(), ctx.PlayerID()))

	return nil

}
