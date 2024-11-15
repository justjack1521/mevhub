package query

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
)

type SearchPlayerQuery struct {
	LobbyID   uuid.UUID
	PartySlot int
}

func NewSearchPlayerQuery(id uuid.UUID, slot int) SearchPlayerQuery {
	return SearchPlayerQuery{LobbyID: id, PartySlot: slot}
}

func (c SearchPlayerQuery) CommandName() string {
	return "search.player"
}

type SearchPlayerQueryHandler struct {
	ParticipantRepository lobby.ParticipantRepository
	SessionRepository     session.InstanceReadRepository
	SummaryRepository     lobby.PlayerSummaryReadRepository
}

func NewSearchPlayerQueryHandler(participant lobby.ParticipantRepository, session session.InstanceReadRepository, summary lobby.PlayerSummaryReadRepository) *SearchPlayerQueryHandler {
	return &SearchPlayerQueryHandler{ParticipantRepository: participant, SessionRepository: session, SummaryRepository: summary}
}

func (h *SearchPlayerQueryHandler) Handle(ctx Context, qry SearchPlayerQuery) (lobby.PlayerSummary, error) {

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, qry.LobbyID, qry.PartySlot)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	current, err := h.SessionRepository.QueryByID(ctx, participant.UserID)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	summary, err := h.SummaryRepository.Query(ctx, current.PlayerID)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	return summary, nil

}
