package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type LobbyReadyCommand struct {
	BasicCommand
}

func (c LobbyReadyCommand) CommandName() string {
	return "lobby.ready"
}

func NewLobbyReadyCommand() *LobbyReadyCommand {
	return &LobbyReadyCommand{}
}

type LobbyReadyCommandHandler struct {
	EventPublisher             *mevent.Publisher
	SessionRepository          session.InstanceReadRepository
	InstanceRepository         port.LobbyInstanceRepository
	QuestRepository            port.QuestRepository
	LobbyPlayerQueueRepository port.MatchLobbyPlayerQueueWriteRepository
}

func NewLobbyReadyCommandHandler(publisher *mevent.Publisher, sessions session.InstanceReadRepository, lobbies port.LobbyInstanceRepository, quests port.QuestRepository, queues port.MatchLobbyPlayerQueueWriteRepository) *LobbyReadyCommandHandler {
	return &LobbyReadyCommandHandler{EventPublisher: publisher, SessionRepository: sessions, InstanceRepository: lobbies, QuestRepository: quests, LobbyPlayerQueueRepository: queues}
}

func (h *LobbyReadyCommandHandler) Handle(ctx Context, cmd *LobbyReadyCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	instance, err := h.InstanceRepository.QueryByID(ctx, current.LobbyID)
	if err != nil {
		return err
	}

	quest, err := h.QuestRepository.QueryByID(instance.QuestID)
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod == game.FulfillMethodMatch {
		if err := h.LobbyPlayerQueueRepository.RemoveLobbyFromQueue(ctx, quest.Tier.GameMode.ModeIdentifier, instance.QuestID, instance.SysID); err != nil {
			return err
		}
	}

	h.EventPublisher.Notify(lobby.NewInstanceReadyEvent(ctx, instance.SysID, instance.QuestID))

	return nil

}
