package query

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type GameSummaryQuery struct {
	InstanceID uuid.UUID
}

func (g GameSummaryQuery) CommandName() string {
	return "game.summary"
}

func NewGameSummaryQuery(id uuid.UUID) GameSummaryQuery {
	return GameSummaryQuery{InstanceID: id}
}

type GameSummaryQueryHandler struct {
	SummaryRepository game.SummaryReadRepository
}

func NewGameSummaryQueryHandler(summary game.SummaryReadRepository) *GameSummaryQueryHandler {
	return &GameSummaryQueryHandler{SummaryRepository: summary}
}

func (h *GameSummaryQueryHandler) Handle(ctx Context, qry GameSummaryQuery) (game.Summary, error) {
	panic(nil)
}
