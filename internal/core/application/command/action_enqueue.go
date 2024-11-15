package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
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
	SessionRepository session.InstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewEnqueueActionCommandHandler(sessions session.InstanceReadRepository, server *server.GameServerHost) *EnqueueActionCommandHandler {
	return &EnqueueActionCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *EnqueueActionCommandHandler) Handle(ctx Context, cmd *EnqueueActionCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var request = &server.GameActionRequest{
		InstanceID: current.LobbyID,
		Action: &game.PlayerEnqueueAction{
			InstanceID: current.PlayerID,
			PlayerID:   current.PlayerID,
			ActionType: cmd.PlayerActionType,
			SlotIndex:  cmd.SlotIndex,
			Target:     cmd.Target,
			ElementID:  cmd.ElementID,
		},
	}

	h.GameServerHost.ActionChannel <- request

	return nil

}
