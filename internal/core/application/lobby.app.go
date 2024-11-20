package application

import (
	"context"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/application/factory"
	"mevhub/internal/core/application/query"
	"mevhub/internal/core/application/service"
	"mevhub/internal/core/application/subscriber"
	"mevhub/internal/core/application/worker"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/decorator"
)

type LobbyApplication struct {
	Queries     *LobbyApplicationQueries
	Commands    *LobbyApplicationCommands
	Translators *LobbyApplicationTranslators
	consumers   []ApplicationConsumer
	subscribers []ApplicationSubscriber
	services    *ApplicationServices
}

type LobbyApplicationQueries struct {
	SearchLobby  SearchLobbyQueryHandler
	SearchPlayer SearchPlayerQueryHandler
}

type LobbyApplicationCommands struct {
	SessionCreate      SessionCreateCommandHandler
	SessionEnd         SessionEndCommandHandler
	LobbyCreate        LobbyCreateCommandHandler
	LobbyCancel        LobbyCancelCommandHandler
	LobbyReady         LobbyReadyCommandHandler
	LobbyStart         LobbyStartCommandHandler
	LobbyStamp         LobbyStampCommandHandler
	ParticipantJoin    ParticipantJoinCommandHandler
	ParticipantReady   ParticipantReadyCommandHandler
	ParticipantUnready ParticipantUnreadyCommandHandler
	ParticipantFind    ParticipantFindCommandHandler
	ParticipantWatch   ParticipantWatchCommandHandler
}

type LobbyApplicationTranslators struct {
	LobbySummary    translate.LobbySummaryTranslator
	LobbyPlayerSlot translate.LobbyPlayerSlotSummaryTranslator
	LobbyPlayer     translate.LobbyPlayerSummaryTranslator
}

func NewLobbyApplication(core *CoreApplication) *LobbyApplication {
	var application = &LobbyApplication{
		consumers: []ApplicationConsumer{
			consumer.NewClientDisconnectConsumer(core.Services.EventPublisher, core.Services.RabbitMQConnection),
		},
	}
	application.Queries = &LobbyApplicationQueries{
		SearchLobby:  application.NewSearchLobbyQueryHandler(core),
		SearchPlayer: application.NewSearchPlayerQueryHandler(core),
	}
	application.Commands = &LobbyApplicationCommands{
		SessionCreate:      application.NewSessionCreateCommandHandler(core),
		SessionEnd:         application.NewSessionEndCommandHandler(core),
		LobbyCreate:        application.NewLobbyCreateCommandHandler(core),
		LobbyCancel:        application.NewLobbyCancelCommandHandler(core),
		LobbyReady:         application.NewLobbyReadyCommandHandler(core),
		LobbyStart:         application.NewLobbyStartCommandHandler(core),
		LobbyStamp:         application.NewLobbyStampCommandHandler(core),
		ParticipantWatch:   application.NewParticipantWatchCommandHandler(core),
		ParticipantJoin:    application.NewParticipantJoinCommandHandler(core),
		ParticipantReady:   application.NewParticipantReadyCommandHandler(core),
		ParticipantUnready: application.NewParticipantUnreadyCommandHandler(core),
		ParticipantFind:    application.NewParticipantFindCommandHandler(core),
	}
	application.Translators = &LobbyApplicationTranslators{
		LobbySummary:    translate.NewLobbySummaryTranslator(),
		LobbyPlayerSlot: translate.NewLobbyPlayerSlotSummaryTranslator(),
		LobbyPlayer:     translate.NewLobbyPlayerSummaryTranslator(),
	}
	application.subscribers = []ApplicationSubscriber{
		subscriber.NewLobbyNotificationChanneler(core.Services.EventPublisher, core.Services.Redis, core.Services.RabbitMQConnection, memory.NewLobbyChannelRepository(core.Services.Redis)),
		subscriber.NewLobbySummaryWriter(core.Services.EventPublisher, core.repositories.Quests, core.data.LobbySummaries),
		subscriber.NewLobbySearchWriter(core.Services.EventPublisher, core.repositories.Quests, core.data.LobbySearch),
		subscriber.NewLobbyQueueWriter(core.Services.EventPublisher, core.data.Lobbies, core.repositories.Quests, core.data.MatchLobbyQueue, core.data.LobbyParticipants),
		subscriber.NewLobbyPlayerQueueWriter(core.Services.EventPublisher, core.data.MatchPlayerQueue, core.repositories.Quests, core.data.LobbyParticipants, core.data.LobbyPlayerSummaries),
		subscriber.NewLobbyChannelEventNotifier(core.Services.EventPublisher, core.data.LobbyPlayerSummaries, application.Translators.LobbyPlayer),
		subscriber.NewLobbyClientNotifier(core.Services.EventPublisher, core.Services.Redis),
		subscriber.NewSessionLobbyWriter(core.Services.EventPublisher, core.data.Sessions),
		subscriber.NewLobbyInstanceWriter(core.Services.EventPublisher, core.data.Lobbies),
		subscriber.NewLobbyParticipantWriter(core.Services.EventPublisher, core.data.LobbyParticipants),
	}

	var lobbyDispatcher = service.NewLobbyMatchmakingDispatcher(core.Services.EventPublisher, core.repositories.Quests, core.data.Lobbies, core.data.Games, factory.NewGameInstanceFactory(core.repositories.Quests))

	var soloLobbyQueueWorker = worker.NewLobbyMatchmakingQueueWorker(context.Background(), game.ModeIdentifierCompSingle, core.data.MatchLobbyQueue, lobbyDispatcher)
	go soloLobbyQueueWorker.Run()

	var duoLobbyQueueWorker = worker.NewLobbyMatchmakingQueueWorker(context.Background(), game.ModeIdentifierCompDuo, core.data.MatchLobbyQueue, lobbyDispatcher)
	go duoLobbyQueueWorker.Run()

	var lobbyPlayerDispatcher = service.NewPlayerMatchmakingDispatcher(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.repositories.Quests, core.data.LobbyParticipants)

	var duoLobbyPlayerQueueWorker = worker.NewLobbyPlayerMatchmakingQueueWorker(context.Background(), game.ModeIdentifierCompDuo, core.data.MatchPlayerQueue, lobbyPlayerDispatcher)
	go duoLobbyPlayerQueueWorker.Run()

	return application
}

type SearchLobbyQueryHandler decorator.QueryHandler[query.Context, query.SearchLobbyQuery, []lobby.Summary]
type SearchPlayerQueryHandler decorator.QueryHandler[query.Context, query.SearchPlayerQuery, lobby.PlayerSummary]

func (a *LobbyApplication) NewSearchLobbyQueryHandler(core *CoreApplication) SearchLobbyQueryHandler {
	var actual = query.NewSearchLobbyQueryHandler(core.data.LobbySearch, service.NewSummaryQueryService(core.data.Lobbies, core.data.LobbyParticipants, core.data.LobbySummaries, core.data.LobbyPlayerSummaries))
	return actual
}

func (a *LobbyApplication) NewSearchPlayerQueryHandler(core *CoreApplication) SearchPlayerQueryHandler {
	var actual = query.NewSearchPlayerQueryHandler(core.data.LobbyParticipants, core.data.Sessions, core.data.LobbyPlayerSummaries)
	return actual
}

type SessionCreateCommandHandler decorator.CommandHandler[command.Context, *command.SessionCreateCommand]
type SessionEndCommandHandler decorator.CommandHandler[command.Context, *command.SessionEndCommand]
type LobbyCreateCommandHandler decorator.CommandHandler[command.Context, *command.LobbyCreateCommand]
type LobbyCancelCommandHandler decorator.CommandHandler[command.Context, *command.LobbyCancelCommand]
type LobbyStartCommandHandler decorator.CommandHandler[command.Context, *command.LobbyStartCommand]
type LobbyReadyCommandHandler decorator.CommandHandler[command.Context, *command.LobbyReadyCommand]
type LobbyStampCommandHandler decorator.CommandHandler[command.Context, *command.LobbyStampCommand]
type ParticipantJoinCommandHandler decorator.CommandHandler[command.Context, *command.ParticipantJoinCommand]
type ParticipantLeaveCommandHandler decorator.CommandHandler[command.Context, *command.ParticipantLeaveCommand]
type ParticipantReadyCommandHandler decorator.CommandHandler[command.Context, *command.ParticipantReadyCommand]
type ParticipantUnreadyCommandHandler decorator.CommandHandler[command.Context, *command.ParticipantUnreadyCommand]
type ParticipantFindCommandHandler decorator.CommandHandler[command.Context, *command.ParticipantFindCommand]
type ParticipantWatchCommandHandler decorator.CommandHandler[command.Context, *command.WatchLobbyCommand]

func (a *LobbyApplication) NewSessionCreateCommandHandler(core *CoreApplication) SessionCreateCommandHandler {
	var actual = command.NewSessionCreateCommandHandler(core.Services.EventPublisher, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.SessionCreateCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewSessionEndCommandHandler(core *CoreApplication) SessionEndCommandHandler {
	var actual = command.NewSessionEndCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.SessionEndCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLobbyCreateCommandHandler(core *CoreApplication) LobbyCreateCommandHandler {
	var actual = command.NewLobbyCreateCommandHandler(core.Services.EventPublisher, core.data.Lobbies, core.repositories.Quests, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LobbyCreateCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLobbyCancelCommandHandler(core *CoreApplication) LobbyCancelCommandHandler {
	var actual = command.NewLobbyCancelCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LobbyCancelCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLobbyStartCommandHandler(core *CoreApplication) LobbyStartCommandHandler {
	var actual = command.NewLobbyStartCommandHandler(core.data.Sessions, core.data.Lobbies, core.data.Games, factory.NewGameInstanceFactory(core.repositories.Quests))
	return decorator.NewStandardCommandDecorator[command.Context, *command.LobbyStartCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantFindCommandHandler(core *CoreApplication) ParticipantFindCommandHandler {
	var actual = command.NewParticipantFindCommandHandler(core.repositories.Quests, core.data.MatchPlayerQueue, core.data.LobbyPlayerSummaries)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ParticipantFindCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLobbyReadyCommandHandler(core *CoreApplication) LobbyReadyCommandHandler {
	var actual = command.NewLobbyReadyCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.repositories.Quests, core.data.MatchPlayerQueue)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LobbyReadyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantJoinCommandHandler(core *CoreApplication) ParticipantJoinCommandHandler {
	var actual = command.NewParticipantJoinCommandHandler(core.Services.EventPublisher, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ParticipantJoinCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantLeaveCommandHandler(core *CoreApplication) ParticipantLeaveCommandHandler {
	var actual = command.NewParticipantLeaveCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ParticipantLeaveCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantReadyCommandHandler(core *CoreApplication) ParticipantReadyCommandHandler {
	var actual = command.NewParticipantReadyCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ParticipantReadyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantUnreadyCommandHandler(core *CoreApplication) ParticipantUnreadyCommandHandler {
	var actual = command.NewParticipantUnreadyCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ParticipantUnreadyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewParticipantWatchCommandHandler(core *CoreApplication) ParticipantWatchCommandHandler {
	var actual = command.NewWatchLobbyCommandHandler(core.Services.EventPublisher)
	return decorator.NewStandardCommandDecorator[command.Context, *command.WatchLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLobbyStampCommandHandler(core *CoreApplication) LobbyStampCommandHandler {
	var actual = command.NewLobbyStampCommandHandler(core.Services.EventPublisher, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LobbyStampCommand](core.Services.EventPublisher, actual)
}
