package query

import (
	uuid "github.com/satori/go.uuid"
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
}
