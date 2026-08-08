package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
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

	instance, err := session.NewInstance(ctx.UserID(), ctx.PlayerID())
	if err != nil {
		return err
	}

	// A relaunching client re-issues SessionCreate while its previous session
	// may still reference a live lobby or game. Creating blind would wipe
	// LobbyID/GameID and orphan the player in the live game (there is no
	// ordering guarantee between this gRPC call and the AMQP ClientConnected
	// event that performs relaunch detection). Carry the previous state over;
	// HandlePlayerConnected remains the sole authority for reconnect/relaunch.
	exists, err := h.SessionRepository.Exists(ctx, ctx.UserID())
	if err != nil {
		return err
	}
	if exists {
		existing, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
		if err != nil {
			return err
		}
		instance.LobbyID = existing.LobbyID
		instance.GameID = existing.GameID
		instance.PartySlot = existing.PartySlot
		instance.DeckIndex = existing.DeckIndex
		instance.DisconnectSessionID = existing.DisconnectSessionID
		instance.LastConnEventAt = existing.LastConnEventAt
		instance.DisconnectedAt = existing.DisconnectedAt
		instance.CurrentSessionID = existing.CurrentSessionID
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
