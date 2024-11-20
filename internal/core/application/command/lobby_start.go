package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/factory"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyStartCommand struct {
	BasicCommand
	GameID uuid.UUID
}

func (c LobbyStartCommand) CommandName() string {
	return "lobby.start"
}

func NewLobbyStartCommand() *LobbyStartCommand {
	return &LobbyStartCommand{
		GameID: uuid.NewV4(),
	}
}

type LobbyStartCommandHandler struct {
	SessionRepository       port.SessionInstanceReadRepository
	LobbyInstanceRepository port.LobbyInstanceRepository
	GameInstanceRepository  port.GameInstanceRepository
	GameInstanceFactory     *factory.GameInstanceFactory
}

func NewLobbyStartCommandHandler(sessions port.SessionInstanceReadRepository, lobbies port.LobbyInstanceRepository, games port.GameInstanceRepository, factory *factory.GameInstanceFactory) *LobbyStartCommandHandler {
	return &LobbyStartCommandHandler{SessionRepository: sessions, LobbyInstanceRepository: lobbies, GameInstanceRepository: games, GameInstanceFactory: factory}
}

func (h *LobbyStartCommandHandler) Handle(ctx Context, cmd *LobbyStartCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	instance, err := h.LobbyInstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		return err
	}

	if err := instance.CanStart(ctx.PlayerID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewInstanceStartedEvent(ctx, instance.SysID))

	result, err := h.GameInstanceFactory.Create(cmd.GameID, instance)
	if err != nil {
		return err
	}

	if err := h.GameInstanceRepository.Create(ctx, result); err != nil {
		return err
	}

	cmd.QueueEvent(game.NewInstanceCreatedEvent(ctx, result.SysID))

	return nil

}
