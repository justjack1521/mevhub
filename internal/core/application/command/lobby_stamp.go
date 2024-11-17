package command

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/domain/session"
)

type LobbyStampCommand struct {
	BasicCommand
	StampID uuid.UUID
}

func NewLobbyStampCommand(id uuid.UUID) *LobbyStampCommand {
	return &LobbyStampCommand{StampID: id}
}

func (c LobbyStampCommand) CommandName() string {
	return "lobby.stamp"
}

type LobbyStampCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository session.InstanceReadRepository
}

func NewLobbyStampCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository) *LobbyStampCommandHandler {
	return &LobbyStampCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *LobbyStampCommandHandler) Handle(ctx Context, cmd *LobbyStampCommand) error {

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

	cmd.QueueEvent(consumer.NewLobbyClientNotificationEvent(ctx, protomulti.MultiLobbyNotificationType_STAMP_SEND, current.LobbyID, bytes))

	return nil

}
