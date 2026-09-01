package command

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/port"
)

type LobbyChatCommand struct {
	BasicCommand
	Message string
}

func NewLobbyChatCommand(message string) *LobbyChatCommand {
	return &LobbyChatCommand{Message: message}
}

func (c LobbyChatCommand) CommandName() string {
	return "lobby.chat"
}

type LobbyChatCommandHandler struct {
	EventPublisher    *mevent.Publisher
	SessionRepository port.SessionInstanceReadRepository
}

func NewLobbyChatCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceReadRepository) *LobbyChatCommandHandler {
	return &LobbyChatCommandHandler{EventPublisher: publisher, SessionRepository: sessions}
}

func (h *LobbyChatCommandHandler) Handle(ctx Context, cmd *LobbyChatCommand) error {

	message, err := sanitiseChatMessage(cmd.Message)
	if err != nil {
		return err
	}

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var notification = &protomulti.LobbyChatNotification{
		LobbyId:   current.LobbyID.String(),
		PartySlot: int32(current.PartySlot),
		Message:   message,
	}

	bytes, err := notification.MarshallBinary()
	if err != nil {
		return err
	}

	cmd.QueueEvent(consumer.NewLobbyClientNotificationEvent(ctx, protomulti.MultiLobbyNotificationType_LOBBY_NOTIFY_CHAT, current.LobbyID, bytes))

	return nil

}
