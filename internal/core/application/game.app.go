package application

import (
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/application/query"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/application/subscriber"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/decorator"
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
	ReadyPlayer   ReadyPlayerCommandHandler
	EnqueueAction EnqueueActionCommandHandler
	DequeueAction DequeueActionCommandHandler
	LockAction    LockActionCommandHandler
}

type GameApplicationTranslators struct {
	PlayerParticipant translate.GamePlayerParticipantTranslator
}

func NewGameApplication(core *CoreApplication) *GameApplication {

	var svr = server.NewGameServerHost(core.Services.RabbitMQConnection, core.Services.Logger, server.NewGameServerFactory([]server.GameServerFactoryBuildAction{
		server.GameServerFactoryLoggingBuildAction(core.Services.Logger),
	}))
	go svr.Run()

	var application = &GameApplication{
		consumers: []ApplicationConsumer{},
	}

	application.Queries = &GameApplicationQueries{
		GameSummary: query.NewGameSummaryQueryHandler(core.data.Sessions, core.data.Games, core.data.GameParticipants),
	}

	application.Commands = &GameApplicationCommands{
		ReadyPlayer:   command.NewReadyPlayerCommandHandler(core.data.Sessions, svr),
		EnqueueAction: command.NewEnqueueActionCommandHandler(core.data.Sessions, svr),
		DequeueAction: command.NewDequeueActionCommandHandler(core.data.Sessions, svr),
		LockAction:    command.NewLockActionCommandHandler(core.data.Sessions, svr),
	}

	application.Translators = &GameApplicationTranslators{
		PlayerParticipant: translate.NewGameParticipantTranslator(),
	}

	application.subscribers = []ApplicationSubscriber{
		subscriber.NewGameChannelEventNotifier(core.Services.EventPublisher),
		subscriber.NewGameInstanceWriter(core.Services.EventPublisher, core.data.Lobbies, game.NewInstanceFactory(core.repositories.Quests), core.data.Games),
		subscriber.NewGameParticipantWriter(core.Services.EventPublisher, core.data.LobbyParticipants, game.NewPlayerParticipantFactory(core.data.GamePlayerLoadouts), core.data.GameParticipants),
		subscriber.NewGameChannelServerWriter(svr, core.Services.EventPublisher, core.data.Games, core.data.GameParticipants),
	}

	return application
}

type ReadyPlayerCommandHandler decorator.CommandHandler[command.Context, *command.ReadyPlayerCommand]
type EnqueueActionCommandHandler decorator.CommandHandler[command.Context, *command.EnqueueActionCommand]
type DequeueActionCommandHandler decorator.CommandHandler[command.Context, *command.DequeueActionCommand]
type LockActionCommandHandler decorator.CommandHandler[command.Context, *command.LockActionCommand]

type GameSummaryQueryHandler decorator.QueryHandler[query.Context, query.GameSummaryQuery, game.Summary]
