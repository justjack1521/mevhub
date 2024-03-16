package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/session"
)

type CreateSessionCommand struct {
	ClientID uuid.UUID
	PlayerID uuid.UUID
}

func NewCreateSessionCommand(client, player uuid.UUID) CreateSessionCommand {
	return CreateSessionCommand{ClientID: client, PlayerID: player}
}

func (c CreateSessionCommand) CommandName() string {
	return "create.session"
}

type CreateSessionCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository session.InstanceWriteRepository
}

func NewCreateSessionCommandHandler(publisher *mevent.Publisher, sessions session.InstanceWriteRepository) *CreateSessionCommandHandler {
	return &CreateSessionCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *CreateSessionCommandHandler) Handle(ctx *Context, cmd CreateSessionCommand) error {

	var instance = &session.Instance{
		ClientID: cmd.ClientID,
		PlayerID: cmd.PlayerID,
	}

	if err := h.SessionRepository.Create(ctx, instance); err != nil {
		return err
	}

	h.EventPublisher.Notify(session.NewInstanceCreatedEvent(ctx, instance.ClientID, instance.PlayerID))

	return nil

}
