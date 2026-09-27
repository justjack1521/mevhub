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
	SearchRepository  port.LobbySearchReadRepository
	SummaryRepository port.LobbySearchSummaryRepository
}

func NewSearchLobbyQueryHandler(lobbies port.LobbySearchReadRepository, summaries port.LobbySearchSummaryRepository) *SearchLobbyQueryHandler {
	return &SearchLobbyQueryHandler{SearchRepository: lobbies, SummaryRepository: summaries}
}

func (h *SearchLobbyQueryHandler) Handle(ctx Context, qry SearchLobbyQuery) ([]lobby.Summary, error) {

	if qry.party != "" {
		summary, err := h.SummaryRepository.QueryByPartyID(ctx, qry.party)
		if err != nil {
			return nil, err
		}
		return []lobby.Summary{summary}, nil
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
		if !hasAvailableSlot(summary) {
			continue
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil

}

// hasAvailableSlot reports whether any empty slot could be joined. Role
// restrictions are not filtered here: the searcher's role isn't known until
// join, and the restriction is returned on each slot for the client to guard.
func hasAvailableSlot(summary lobby.Summary) bool {
	for _, slot := range summary.Players {
		if !uuid.Equal(slot.PlayerSummary.Identity.PlayerID, uuid.Nil) {
			continue
		}
		if slot.Locked || slot.BotControl {
			continue
		}
		return true
	}
	return false
}
