package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
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
	SessionReadRepository  session.InstanceReadRepository
	SessionWriteRepository session.InstanceWriteRepository
}

func NewSessionEndCommandHandler(publisher *mevent.Publisher, read session.InstanceReadRepository, write session.InstanceWriteRepository) *SessionEndCommandHandler {
	return &SessionEndCommandHandler{EventPublisher: publisher, SessionReadRepository: read, SessionWriteRepository: write}
}

func (h *SessionEndCommandHandler) Handle(ctx Context, cmd *SessionEndCommand) error {

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