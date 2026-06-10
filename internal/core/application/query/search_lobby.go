package query

import (
	uuid "github.com/satori/go.uuid"

	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type SearchLobbyQuery struct {
	party string
	query lobby.SearchQuery
}

func (s SearchLobbyQuery) CommandName() string {
	return "search.lobby"
}

func NewSearchLobbyQuery(qry lobby.SearchQuery, party string) SearchLobbyQuery {
	return SearchLobbyQuery{query: qry, party: party}
}

type SearchLobbyQueryHandler struct {
	SearchRepository        port.LobbySearchReadRepository
	SummaryRepository       port.LobbySearchSummaryRepository
	PlayerSummaryRepository port.LobbyPlayerSummaryReadRepository
}

func NewSearchLobbyQueryHandler(lobbies port.LobbySearchReadRepository, summaries port.LobbySearchSummaryRepository, players port.LobbyPlayerSummaryReadRepository) *SearchLobbyQueryHandler {
	return &SearchLobbyQueryHandler{SearchRepository: lobbies, SummaryRepository: summaries, PlayerSummaryRepository: players}
}

func (h *SearchLobbyQueryHandler) Handle(ctx Context, qry SearchLobbyQuery) ([]lobby.Summary, error) {

	if qry.party != "" {
		summary, err := h.SummaryRepository.QueryByPartyID(ctx, qry.party)
		if err != nil {
			return nil, err
		}
		return []lobby.Summary{summary}, nil
	}

	playerSummary, err := h.PlayerSummaryRepository.Query(ctx, ctx.PlayerID())
	var playerRole uuid.UUID
	if err == nil {
		playerRole = playerSummary.Loadout.JobCard.JobCardID
	}

	lobbies, err := h.SearchRepository.Query(ctx, qry.query)
	if err != nil {
		return nil, err
	}

	var summaries = make([]lobby.Summary, 0, len(lobbies))

	for _, value := range lobbies {
		summary, err := h.SummaryRepository.Query(ctx, value.LobbyID)
		if err != nil {
			continue
		}
		if !hasAvailableSlot(summary, playerRole) {
			continue
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil

}

func hasAvailableSlot(summary lobby.Summary, playerRole uuid.UUID) bool {
	for _, slot := range summary.Players {
		if !uuid.Equal(slot.PlayerSummary.Identity.PlayerID, uuid.Nil) {
			continue
		}
		if uuid.Equal(slot.RoleRestriction, uuid.Nil) {
			return true
		}
		if uuid.Equal(slot.RoleRestriction, playerRole) {
			return true
		}
	}
	return false
}
