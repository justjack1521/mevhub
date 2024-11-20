package query

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/factory"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GameSummaryQuery struct {
	GameID uuid.UUID
}

func NewGameSummaryQuery(id uuid.UUID) GameSummaryQuery {
	return GameSummaryQuery{
		GameID: id,
	}
}

func (g GameSummaryQuery) CommandName() string {
	return "query.game.summary"
}

type GameSummaryQueryHandler struct {
	GameInstanceRepository    port.GameInstanceReadRepository
	GamePartyRepository       port.GamePartyReadRepository
	GameParticipantRepository port.GameParticipantReadRepository
	GamePlayerFactory         *factory.GamePlayerFactory
}

func NewGameSummaryQueryHandler(games port.GameInstanceReadRepository, parties port.GamePartyReadRepository, participants port.GameParticipantReadRepository, factory *factory.GamePlayerFactory) *GameSummaryQueryHandler {
	return &GameSummaryQueryHandler{GameInstanceRepository: games, GamePartyRepository: parties, GameParticipantRepository: participants, GamePlayerFactory: factory}
}

func (h *GameSummaryQueryHandler) Handle(ctx Context, cmd GameSummaryQuery) (game.Summary, error) {

	instance, err := h.GameInstanceRepository.Get(ctx, cmd.GameID)
	if err != nil {
		return game.Summary{}, err
	}

	var summary = game.Summary{
		SysID:   instance.SysID,
		Seed:    instance.Seed,
		Parties: make([]game.PartySummary, 0),
	}

	parties, err := h.GamePartyRepository.QueryAll(ctx, cmd.GameID)
	if err != nil {
		return game.Summary{}, err
	}

	for _, value := range parties {

		var party = game.PartySummary{
			SysID:     value.SysID,
			PartyID:   value.PartyID,
			Index:     value.Index,
			PartyName: value.PartyName,
			Players:   make([]game.Player, 0),
		}

		participants, err := h.GameParticipantRepository.QueryAll(ctx, party.SysID)
		if err != nil {
			return game.Summary{}, err
		}

		for _, participant := range participants {
			player, err := h.GamePlayerFactory.Create(ctx, participant)
			if err != nil {
				return game.Summary{}, err
			}
			party.Players = append(party.Players, player)
		}

		summary.Parties = append(summary.Parties, party)
	}

	return summary, nil

}
