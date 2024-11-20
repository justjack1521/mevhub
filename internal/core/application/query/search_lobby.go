package query

import (
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
	InstanceRepository port.LobbyInstanceRepository
	SearchRepository   port.LobbySearchReadRepository
	SummaryRepository  port.LobbySummaryReadRepository
}

func NewSearchLobbyQueryHandler(lobbies port.LobbySearchReadRepository, summaries port.LobbySummaryReadRepository) *SearchLobbyQueryHandler {
	return &SearchLobbyQueryHandler{SearchRepository: lobbies, SummaryRepository: summaries}
}

func (h *SearchLobbyQueryHandler) Handle(ctx Context, qry SearchLobbyQuery) ([]lobby.Summary, error) {

	var summaries = make([]lobby.Summary, 0)

	if qry.party != "" {
		instance, err := h.InstanceRepository.QueryByPartyID(ctx, qry.party)
		if err != nil {
			return nil, err
		}
		summary, err := h.SummaryRepository.Query(ctx, instance.SysID)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	} else {
		lobbies, err := h.SearchRepository.Query(ctx, qry.query)
		if err != nil {
			return nil, err
		}
		for _, value := range lobbies {
			summary, err := h.SummaryRepository.Query(ctx, value.LobbyID)
			if err != nil {
				continue
			}
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil

}
