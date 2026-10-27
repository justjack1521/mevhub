package command

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
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
	EventPublisher     *mevent.Publisher
	SessionRepository  port.SessionInstanceReadRepository
	CooldownRepository port.LobbyStampCooldownRepository
}

func NewLobbyStampCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceReadRepository, cooldowns port.LobbyStampCooldownRepository) *LobbyStampCommandHandler {
	return &LobbyStampCommandHandler{EventPublisher: publisher, SessionRepository: sessions, CooldownRepository: cooldowns}
}

func (h *LobbyStampCommandHandler) Handle(ctx Context, cmd *LobbyStampCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	// A stamp inside the cooldown is dropped rather than refused, as in game: the
	// client already holds the sender back, so this only meets a modified one.
	claimed, err := h.CooldownRepository.Claim(ctx, ctx.UserID(), lobby.StampCooldown)
	if err != nil {
		return err
	}
	if claimed == false {
		return nil
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

	cmd.QueueEvent(consumer.NewLobbyClientNotificationEvent(ctx, protomulti.MultiLobbyNotificationType_LOBBY_NOTIFY_STAMP_SEND, current.LobbyID, bytes))

	return nil

}
