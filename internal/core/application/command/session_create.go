package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type SessionCreateCommand struct {
	BasicCommand
}

func NewSessionCreateCommand() *SessionCreateCommand {
	return &SessionCreateCommand{}
}

func (c SessionCreateCommand) CommandName() string {
	return "session.create"
}

type SessionCreateCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceWriteRepository
}

func NewSessionCreateCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceWriteRepository) *SessionCreateCommandHandler {
	return &SessionCreateCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *SessionCreateCommandHandler) Handle(ctx Context, cmd *SessionCreateCommand) error {

	instance, err := session.NewInstance(ctx.UserID(), ctx.PlayerID())
	if err != nil {
		return err
	}

	if err := h.SessionRepository.Create(ctx, instance); err != nil {
		return err
	}

	cmd.QueueEvent(session.NewInstanceCreatedEvent(ctx, instance.UserID, instance.PlayerID))

	return nil

}
