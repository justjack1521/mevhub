package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
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
	EventPublisher         *mevent.Publisher
	SessionReadRepository  port.SessionInstanceReadRepository
	SessionWriteRepository port.SessionInstanceWriteRepository
}

func NewSessionEndCommandHandler(publisher *mevent.Publisher, read port.SessionInstanceReadRepository, write port.SessionInstanceWriteRepository) *SessionEndCommandHandler {
	return &SessionEndCommandHandler{EventPublisher: publisher, SessionReadRepository: read, SessionWriteRepository: write}
}

func (h *SessionEndCommandHandler) Handle(ctx Context, cmd *SessionEndCommand) error {

	exists, err := h.SessionReadRepository.Exists(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if exists == false {
		return nil
	}

	instance, err := h.SessionReadRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if err := h.SessionWriteRepository.Delete(ctx, instance); err != nil {
		return err
	}

	cmd.QueueEvent(session.NewInstanceDeletedEvent(ctx, instance.UserID, instance.UserID, instance.PlayerID))

	return nil

}
