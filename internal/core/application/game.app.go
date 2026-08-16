package application

import (
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/application/factory"
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
	ReadyPlayer       ReadyPlayerCommandHandler
	EnqueueAction     EnqueueActionCommandHandler
	DequeueAction     DequeueActionCommandHandler
	LockAction        LockActionCommandHandler
	SubmitHPConsensus SubmitHPConsensusCommandHandler
	GameCatchUp       GameCatchUpCommandHandler
}

type GameApplicationTranslators struct {
	Summary translate.GameSummaryTranslator
	Player  translate.GamePlayerTranslator
}

func NewGameApplication(core *CoreApplication) *GameApplication {

	var svr = server.NewGameServerHost(core.Services.Logger, server.NewGameServerFactory([]server.GameServerFactoryBuildAction{
		server.GameServerFactoryPublisherBuildAction(server.NewGameServerRabbitMQNotifier(core.Services.RabbitMQConnection), core.Services.EventPublisher),
		server.GameServerFactoryLoggingBuildAction(core.Services.Logger),
	}), core.Services.EventPublisher)
	go svr.Run()

	var application = &GameApplication{
		consumers: []ApplicationConsumer{},
	}

	application.Queries = &GameApplicationQueries{
		GameSummary: query.NewGameSummaryQueryHandler(core.data.Games, core.data.GameParties, core.data.GameParticipants, factory.NewGamePlayerFactory(core.data.GamePlayerLoadouts)),
	}

	application.Commands = &GameApplicationCommands{
		ReadyPlayer:       command.NewReadyPlayerCommandHandler(core.data.Sessions, svr),
		EnqueueAction:     command.NewEnqueueActionCommandHandler(core.data.Sessions, svr),
		DequeueAction:     command.NewDequeueActionCommandHandler(core.data.Sessions, svr),
		LockAction:        command.NewLockActionCommandHandler(core.data.Sessions, svr),
		SubmitHPConsensus: command.NewSubmitHPConsensusCommandHandler(core.data.Sessions, svr),
		GameCatchUp:       command.NewGameCatchUpCommandHandler(core.data.Sessions, svr),
	}

	application.Translators = &GameApplicationTranslators{
		Summary: translate.NewGameSummaryTranslator(),
		Player:  translate.NewGamePlayerTranslator(),
	}

	// Order matters: the publisher notifies handlers in subscription order, and
	// all three of these subscribe to game.InstanceCreatedEvent's fan-out. The
	// party and participant writers must have written their read models before
	// the channel writer queues the game server for registration, since the
	// registration callback populates the live game from those read models.
	application.subscribers = []ApplicationSubscriber{
		subscriber.NewGamePartyWriter(core.Services.EventPublisher, core.data.Games, core.data.LobbySummaries, core.data.GameParties, core.Services.Logger),
		subscriber.NewGameParticipantWriter(core.Services.EventPublisher, core.data.LobbyParticipants, core.data.GameParticipants, core.Services.Logger),
		subscriber.NewGameChannelServerWriter(svr, core.Services.EventPublisher, core.data.Games, core.data.GameParties, core.data.GameParticipants, core.Services.Logger),
		subscriber.NewGameLoadoutEvictionSubscriber(core.Services.EventPublisher, core.data.GamePlayerLoadouts),
	}

	return application
}

type ReadyPlayerCommandHandler decorator.CommandHandler[command.Context, *command.ReadyPlayerCommand]
type EnqueueActionCommandHandler decorator.CommandHandler[command.Context, *command.EnqueueActionCommand]
type DequeueActionCommandHandler decorator.CommandHandler[command.Context, *command.DequeueActionCommand]
type LockActionCommandHandler decorator.CommandHandler[command.Context, *command.LockActionCommand]
type SubmitHPConsensusCommandHandler decorator.CommandHandler[command.Context, *command.SubmitHPConsensusCommand]
type GameCatchUpCommandHandler decorator.CommandHandler[command.Context, *command.GameCatchUpCommand]

type GameSummaryQueryHandler decorator.QueryHandler[query.Context, query.GameSummaryQuery, game.Summary]
