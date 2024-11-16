package command

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
)

type FindLobbyCommand struct {
	BasicCommand
	QuestID    uuid.UUID
	DeckIndex  int
	UseStamina bool
}

func (c FindLobbyCommand) CommandName() string {
	return "find.lobby"
}

type FindLobbyCommandHandler struct {
	QuestRepository       port.QuestRepository
	MatchmakingRepository port.MatchPlayerQueueWriteRepository
}

func (h *FindLobbyCommandHandler) Handle(ctx Context, cmd *FindLobbyCommand) error {

	quest, err := h.QuestRepository.QueryByID(cmd.QuestID)
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return errors.New("cannot find lobbies for this game mode")
	}

	var entry = match.PlayerQueueEntry{
		UserID:    ctx.UserID(),
		QuestID:   cmd.QuestID,
		DeckLevel: cmd.DeckIndex,
		JoinedAt:  time.Now().UTC(),
	}

	if err := h.MatchmakingRepository.AddPlayerToQueue(ctx, quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
		return err
	}

	return nil

}
