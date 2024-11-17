package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
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
	SessionRepository session.InstanceWriteRepository
}

func NewSessionCreateCommandHandler(publisher *mevent.Publisher, sessions session.InstanceWriteRepository) *SessionCreateCommandHandler {
	return &SessionCreateCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *SessionCreateCommandHandler) Handle(ctx Context, cmd *SessionCreateCommand) error {

	var instance = &session.Instance{
		UserID:   ctx.UserID(),
		PlayerID: ctx.PlayerID(),
	}

	if err := h.SessionRepository.Create(ctx, instance); err != nil {
		return err
	}

	cmd.QueueEvent(session.NewInstanceCreatedEvent(ctx, instance.UserID, instance.PlayerID))

	return nil

}
