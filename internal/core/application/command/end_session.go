package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
)

type EndSessionCommand struct {
	BasicCommand
}

func (e EndSessionCommand) CommandName() string {
	return "session.end"
}

func NewEndSessionCommand() *EndSessionCommand {
	return &EndSessionCommand{}
}

type EndSessionCommandHandler struct {
	EventPublisher         *mevent.Publisher
	SessionReadRepository  session.InstanceReadRepository
	SessionWriteRepository session.InstanceWriteRepository
}

func NewEndSessionCommandHandler(publisher *mevent.Publisher, read session.InstanceReadRepository, write session.InstanceWriteRepository) *EndSessionCommandHandler {
	return &EndSessionCommandHandler{EventPublisher: publisher, SessionReadRepository: read, SessionWriteRepository: write}
}

func (h *EndSessionCommandHandler) Handle(ctx Context, cmd *EndSessionCommand) error {

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
