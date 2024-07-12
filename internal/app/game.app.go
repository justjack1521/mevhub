package app

import (
	"mevhub/internal/app/query"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/game"
)

type GameApplication struct {
	Queries     *GameApplicationQueries
	Commands    *GameApplicationCommands
	consumers   []ApplicationConsumer
	subscribers []ApplicationSubscriber
}

type GameApplicationQueries struct {
	GameSummary GameSummaryQueryHandler
}

type GameApplicationCommands struct {
}

type GameApplicationTranslators struct {
}

func NewGameApplication(core *CoreApplication) *GameApplication {
	var application = &GameApplication{
		consumers:   []ApplicationConsumer{},
		subscribers: []ApplicationSubscriber{},
	}
	application.Queries = &GameApplicationQueries{
		GameSummary: application.NewGameSummaryQueryHandler(core),
	}
	application.Commands = &GameApplicationCommands{}
	return application
}

type GameSummaryQueryHandler decorator.QueryHandler[query.Context, query.GameSummaryQuery, game.Summary]

func (a *GameApplication) NewGameSummaryQueryHandler(core *CoreApplication) GameSummaryQueryHandler {
	var actual = query.NewGameSummaryQueryHandler(core.data.GameSummary)
	return actual
}
