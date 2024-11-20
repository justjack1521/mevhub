package query

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
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
}

func NewGameSummaryQueryHandler() *GameSummaryQueryHandler {
	return &GameSummaryQueryHandler{}
}

func (h *GameSummaryQueryHandler) Handle(ctx Context, cmd GameSummaryQuery) (game.Summary, error) {

	panic(nil)

}
