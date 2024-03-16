package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type WatchLobbyCommand struct {
	LobbyID uuid.UUID
}

func (c WatchLobbyCommand) CommandName() string {
	return "inspect.lobby"
}

func NewWatchLobbyCommand(lobby uuid.UUID) WatchLobbyCommand {
	return WatchLobbyCommand{LobbyID: lobby}
}

type WatchLobbyCommandHandler struct {
	EventPublisher *mevent.Publisher
}

func NewWatchLobbyCommandHandler(publisher *mevent.Publisher) *WatchLobbyCommandHandler {
	return &WatchLobbyCommandHandler{EventPublisher: publisher}
}

func (h *WatchLobbyCommandHandler) Handle(ctx *Context, cmd WatchLobbyCommand) error {

	h.EventPublisher.Notify(lobby.NewWatcherAddedEvent(ctx, cmd.LobbyID, ctx.ClientID, ctx.Session.PlayerID))

	return nil

}
