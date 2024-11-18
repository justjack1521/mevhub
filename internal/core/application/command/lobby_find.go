package command

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
)

type ParticipantFindCommand struct {
	BasicCommand
	QuestID    uuid.UUID
	DeckIndex  int
	UseStamina bool
}

func NewParticipantFindCommand(quest uuid.UUID, index int, stamina bool) *ParticipantFindCommand {
	return &ParticipantFindCommand{QuestID: quest, DeckIndex: index, UseStamina: stamina}
}

func (c ParticipantFindCommand) CommandName() string {
	return "participant.find"
}

type ParticipantFindCommandHandler struct {
	QuestRepository         port.QuestRepository
	MatchmakingRepository   port.MatchLobbyPlayerQueueWriteRepository
	PlayerSummaryRepository port.LobbyPlayerSummaryReadRepository
}

func NewParticipantFindCommandHandler(quests port.QuestRepository, queue port.MatchLobbyPlayerQueueWriteRepository, players port.LobbyPlayerSummaryReadRepository) *ParticipantFindCommandHandler {
	return &ParticipantFindCommandHandler{QuestRepository: quests, MatchmakingRepository: queue, PlayerSummaryRepository: players}
}

func (h *ParticipantFindCommandHandler) Handle(ctx Context, cmd *ParticipantFindCommand) error {

	quest, err := h.QuestRepository.QueryByID(cmd.QuestID)
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return errors.New("cannot find lobbies for this game mode")
	}

	player, err := h.PlayerSummaryRepository.Query(ctx, ctx.PlayerID())
	if err != nil {
		return err
	}

	var entry = match.PlayerQueueEntry{
		UserID:    ctx.UserID(),
		QuestID:   cmd.QuestID,
		DeckLevel: player.Loadout.CalculateDeckLevel(),
		JoinedAt:  time.Now().UTC(),
	}

	if err := h.MatchmakingRepository.AddPlayerToQueue(ctx, quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
		return err
	}

	return nil

}
