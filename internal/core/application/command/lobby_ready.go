package command

import (
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
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
	SessionRepository          port.SessionInstanceReadRepository
	InstanceRepository         port.LobbyInstanceRepository
	QuestRepository            port.QuestRepository
	LobbyPlayerQueueRepository port.MatchLobbyPlayerQueueWriteRepository
	LobbyQueueRepository       port.MatchLobbyQueueWriteRepository
	ParticipantRepository      port.LobbyParticipantReadRepository
	PlayerSummaryRepository    port.LobbyPlayerSummaryReadRepository
}

func NewLobbyReadyCommandHandler(
	sessions port.SessionInstanceReadRepository,
	lobbies port.LobbyInstanceRepository,
	quests port.QuestRepository,
	playerQueue port.MatchLobbyPlayerQueueWriteRepository,
	lobbyQueue port.MatchLobbyQueueWriteRepository,
	participants port.LobbyParticipantReadRepository,
	players port.LobbyPlayerSummaryReadRepository,
) *LobbyReadyCommandHandler {
	return &LobbyReadyCommandHandler{
		SessionRepository:          sessions,
		InstanceRepository:         lobbies,
		QuestRepository:            quests,
		LobbyPlayerQueueRepository: playerQueue,
		LobbyQueueRepository:       lobbyQueue,
		ParticipantRepository:      participants,
		PlayerSummaryRepository:    players,
	}
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

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return nil
	}

	if err := h.LobbyPlayerQueueRepository.RemoveLobbyFromQueue(ctx, quest.Tier.GameMode.ModeIdentifier, instance.QuestID, instance.SysID); err != nil {
		return err
	}

	participants, err := h.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return err
	}

	var sum, count int
	for _, participant := range participants {
		if !participant.HasPlayer() {
			continue
		}
		summary, err := h.PlayerSummaryRepository.Query(ctx, participant.PlayerID)
		if err != nil {
			return err
		}
		sum += summary.Loadout.CalculateDeckLevel()
		count++
	}

	var average int
	if count > 0 {
		average = sum / count
	}

	var entry = match.LobbyQueueEntry{
		LobbyID:      instance.SysID,
		QuestID:      instance.QuestID,
		AverageLevel: average,
		JoinedAt:     time.Now().UTC(),
	}

	if err := h.LobbyQueueRepository.AddLobbyToQueue(ctx, quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
		return err
	}

	return nil

}
