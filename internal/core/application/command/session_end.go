package command

import (
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
)

type SessionEndCommand struct {
	BasicCommand
}

func (e SessionEndCommand) CommandName() string {
	return "session.end"
}

func NewSessionEndCommand() *SessionEndCommand {
	return &SessionEndCommand{}
}

type SessionEndCommandHandler struct {
	EventPublisher *mevent.Publisher
	read           port.SessionInstanceReadRepository
	write          port.SessionInstanceWriteRepository
}

func NewSessionEndCommandHandler(publisher *mevent.Publisher, read port.SessionInstanceReadRepository, write port.SessionInstanceWriteRepository) *SessionEndCommandHandler {
	return &SessionEndCommandHandler{EventPublisher: publisher, read: read, write: write}
}

func (h *SessionEndCommandHandler) Handle(ctx Context, cmd *SessionEndCommand) error {

	exists, err := h.read.Exists(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if exists == false {
		return nil
	}

	instance, err := h.read.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if err := h.write.Delete(ctx, instance); err != nil {
		return err
	}

	cmd.QueueEvent(session.NewInstanceDeletedEvent(ctx, instance.UserID, instance.PlayerID, instance.LobbyID, instance.GameID, instance.DeckIndex))

	return nil

}
