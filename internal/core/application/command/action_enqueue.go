package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type EnqueueActionCommand struct {
	BasicCommand
	PlayerActionType game.PlayerActionType
	Target           int
	SlotIndex        int
	ElementID        uuid.UUID
}

func (e *EnqueueActionCommand) CommandName() string {
	return "action.enqueue"
}

func NewEnqueueActionCommand(action game.PlayerActionType, target, slot int, element uuid.UUID) *EnqueueActionCommand {
	return &EnqueueActionCommand{
		PlayerActionType: action,
		Target:           target,
		SlotIndex:        slot,
		ElementID:        element,
	}
}

type EnqueueActionCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewEnqueueActionCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *EnqueueActionCommandHandler {
	return &EnqueueActionCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *EnqueueActionCommandHandler) Handle(ctx Context, cmd *EnqueueActionCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  game.NewPlayerEnqueueAction(current.GameID, current.LobbyID, current.PlayerID, cmd.Target, cmd.PlayerActionType, cmd.SlotIndex, cmd.ElementID),
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
