package command

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/app/consumer"
	"mevhub/internal/domain/session"
)

type SendStampCommand struct {
	StampID uuid.UUID
}

func NewSendStampCommand(id uuid.UUID) SendStampCommand {
	return SendStampCommand{StampID: id}
}

func (c SendStampCommand) CommandName() string {
	return "send.stamp"
}

type SendStampCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository session.InstanceReadRepository
}

func NewSendStampCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository) *SendStampCommandHandler {
	return &SendStampCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *SendStampCommandHandler) Handle(ctx Context, cmd SendStampCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var notification = &protomulti.StampSendNotification{
		LobbyId:   current.LobbyID.String(),
		StampId:   cmd.StampID.String(),
		PartySlot: int32(current.PartySlot),
	}

	bytes, err := notification.MarshallBinary()
	if err != nil {
		return err
	}

	var message = consumer.NewLobbyClientNotificationEvent(ctx, protomulti.MultiNotificationType_STAMP_SEND, current.LobbyID, bytes)
	h.EventPublisher.Notify(message)

	return nil

}
