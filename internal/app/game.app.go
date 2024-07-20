package app

import (
	"mevhub/internal/adapter/translate"
	"mevhub/internal/app/query"
	"mevhub/internal/app/server"
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
		GameSummary: query.NewGameSummaryQueryHandler(core.data.Sessions, core.data.Games, core.data.GameParticipants),
	}
	application.Commands = &GameApplicationCommands{}
	application.Translators = &GameApplicationTranslators{
		PlayerParticipant: translate.NewGameParticipantTranslator(),
	}

	var svr = server.NewGameChannelServer(core.Services.RabbitMQConnection, core.Services.Logger)

	application.subscribers = []ApplicationSubscriber{
		subscriber.NewGameChannelEventNotifier(core.Services.EventPublisher),
		subscriber.NewGameInstanceWriter(core.Services.EventPublisher, core.data.Lobbies, game.NewInstanceFactory(core.repositories.Quests), core.data.Games),
		subscriber.NewGameParticipantWriter(core.Services.EventPublisher, core.data.LobbyParticipants, game.NewPlayerParticipantFactory(core.data.GamePlayerLoadouts), core.data.GameParticipants),
		subscriber.NewGameChannelServerWriter(svr, core.Services.EventPublisher, core.data.Games, core.data.GameParticipants),
	}

	return application
}

type GameSummaryQueryHandler decorator.QueryHandler[query.Context, query.GameSummaryQuery, game.Summary]
