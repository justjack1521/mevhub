package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/session"
)

type EndSessionCommand struct {
	SysID    uuid.UUID
	PlayerID uuid.UUID
}

func (e EndSessionCommand) CommandName() string {
	return "end.session"
}

func NewEndSessionCommand(id uuid.UUID, player uuid.UUID) EndSessionCommand {
	return EndSessionCommand{SysID: id, PlayerID: player}
}

type EndSessionCommandHandler struct {
	EventPublisher         *mevent.Publisher
	SessionReadRepository  session.InstanceReadRepository
	SessionWriteRepository session.InstanceWriteRepository
}

func NewEndSessionCommandHandler(publisher *mevent.Publisher, read session.InstanceReadRepository, write session.InstanceWriteRepository) *EndSessionCommandHandler {
	return &EndSessionCommandHandler{EventPublisher: publisher, SessionReadRepository: read, SessionWriteRepository: write}
}

func (h *EndSessionCommandHandler) Handle(ctx *Context, cmd EndSessionCommand) error {

	instance, err := h.SessionReadRepository.QueryByID(ctx, cmd.SysID)
	if err != nil {
		return err
	}

	if err := h.SessionWriteRepository.Delete(ctx, instance); err != nil {
		return err
	}

	h.EventPublisher.Notify(session.NewInstanceDeletedEvent(ctx, instance.ClientID, instance.ClientID, instance.PlayerID))

	return nil

}
