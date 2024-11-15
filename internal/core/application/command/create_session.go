package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
)

type CreateSessionCommand struct {
	BasicCommand
}

func NewCreateSessionCommand() *CreateSessionCommand {
	return &CreateSessionCommand{}
}

func (c CreateSessionCommand) CommandName() string {
	return "session.create"
}

type CreateSessionCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository session.InstanceWriteRepository
}

func NewCreateSessionCommandHandler(publisher *mevent.Publisher, sessions session.InstanceWriteRepository) *CreateSessionCommandHandler {
	return &CreateSessionCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *CreateSessionCommandHandler) Handle(ctx Context, cmd *CreateSessionCommand) error {

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
