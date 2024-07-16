package app

import (
	"mevhub/internal/adapter/translate"
	"mevhub/internal/app/query"
	"mevhub/internal/app/subscriber"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/game"
)

type GameApplication struct {
	Queries     *GameApplicationQueries
	Commands    *GameApplicationCommands
	Translators *GameApplicationTranslators
	consumers   []ApplicationConsumer
	subscribers []ApplicationSubscriber
}

type GameApplicationQueries struct {
	GameSummary GameSummaryQueryHandler
}

type GameApplicationCommands struct {
}

type GameApplicationTranslators struct {
	PlayerParticipant translate.GamePlayerParticipantTranslator
}

func NewGameApplication(core *CoreApplication) *GameApplication {
	var application = &GameApplication{
		consumers: []ApplicationConsumer{},
	}

	application.Queries = &GameApplicationQueries{
		GameSummary: query.NewGameSummaryQueryHandler(core.data.GameInstance, core.data.GamePlayerParticipant),
	}
	application.Commands = &GameApplicationCommands{}
	application.Translators = &GameApplicationTranslators{
		PlayerParticipant: translate.NewGameParticipantTranslator(),
	}

	application.subscribers = []ApplicationSubscriber{
		subscriber.NewGameChannelEventNotifier(core.Services.EventPublisher),
		subscriber.NewGameInstanceWriter(core.Services.EventPublisher, core.data.LobbyInstance, game.NewInstanceFactory(core.repositories.Quests), core.data.GameInstance),
		subscriber.NewGameParticipantWriter(core.Services.EventPublisher, core.data.LobbyParticipant, game.NewPlayerParticipantFactory(core.data.GamePlayerLoadout), core.data.GamePlayerParticipant),
	}

	return application
}

type GameSummaryQueryHandler decorator.QueryHandler[query.Context, query.GameSummaryQuery, game.Summary]
