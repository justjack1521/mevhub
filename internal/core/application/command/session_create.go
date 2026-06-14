package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
	uuid "github.com/satori/go.uuid"
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
	EventPublisher          *mevent.Publisher
	SessionRepository       port.SessionInstanceRepository
	PlayerSummaryRepository port.LobbyPlayerSummaryWriteRepository
}

func NewSessionCreateCommandHandler(publisher *mevent.Publisher, sessions port.SessionInstanceRepository, players port.LobbyPlayerSummaryWriteRepository) *SessionCreateCommandHandler {
	return &SessionCreateCommandHandler{EventPublisher: publisher, SessionRepository: sessions, PlayerSummaryRepository: players}
}

func (h *SessionCreateCommandHandler) Handle(ctx Context, cmd *SessionCreateCommand) error {

	exists, err := h.SessionRepository.Exists(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	if exists {
		existing, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
		if err != nil {
			return err
		}
		if !uuid.Equal(existing.GameID, uuid.Nil) {
			if err := h.SessionRepository.Update(ctx, existing); err != nil {
				return err
			}
			if err := h.PlayerSummaryRepository.Delete(ctx, ctx.PlayerID()); err != nil {
				return err
			}
			cmd.QueueEvent(game.NewPlayerReconnectedEvent(ctx, existing.GameID, existing.UserID, existing.PlayerID))
			return nil
		}
	}

	instance, err := session.NewInstance(ctx.UserID(), ctx.PlayerID())
	if err != nil {
		return err
	}

	if err := h.SessionRepository.Create(ctx, instance); err != nil {
		return err
	}

	if err := h.PlayerSummaryRepository.Delete(ctx, ctx.PlayerID()); err != nil {
		return err
	}

	cmd.QueueEvent(session.NewInstanceCreatedEvent(ctx, instance.UserID, instance.PlayerID))

	return nil

}
