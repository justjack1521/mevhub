package app

import (
	"mevhub/internal/adapter/memory"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/app/command"
	"mevhub/internal/app/query"
	"mevhub/internal/app/subscriber"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/lobby"
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
	CreateSession CreateSessionCommandHandler
	EndSession    EndSessionCommandHandler
	CreateLobby   CreateLobbyCommandHandler
	WatchLobby    WatchLobbyCommandHandler
	JoinLobby     JoinLobbyCommandHandler
	CancelLobby   CancelLobbyCommandHandler
	ReadyLobby    ReadyLobbyCommandHandler
	UnreadyLobby  UnreadyLobbyCommandHandler
	SendStamp     SendStampCommandHandler
}

type LobbyApplicationTranslators struct {
	LobbySummary    translate.LobbySummaryTranslator
	LobbyPlayerSlot translate.LobbyPlayerSlotSummaryTranslator
	LobbyPlayer     translate.LobbyPlayerSummaryTranslator
}

func NewLobbyApplication(core *CoreApplication) *LobbyApplication {
	var application = &LobbyApplication{
		consumers: []ApplicationConsumer{},
	}
	application.Queries = &LobbyApplicationQueries{
		SearchLobby:  application.NewSearchLobbyQueryHandler(core),
		SearchPlayer: application.NewSearchPlayerQueryHandler(core),
	}
	application.Commands = &LobbyApplicationCommands{
		CreateSession: application.NewCreateSessionCommandHandler(core),
		EndSession:    application.NewEndSessionCommandHandler(core),
		CreateLobby:   application.NewCreateLobbyCommandHandler(core),
		CancelLobby:   application.NewCancelLobbyCommandHandler(core),
		WatchLobby:    application.NewWatchLobbyCommandHandler(core),
		JoinLobby:     application.NewJoinLobbyCommandHandler(core),
		ReadyLobby:    application.NewReadyLobbyCommandHandler(core),
		UnreadyLobby:  application.NewUnreadyLobbyCommandHandler(core),
		SendStamp:     application.NewSendStampCommandHandler(core),
	}
	application.Translators = &LobbyApplicationTranslators{
		LobbySummary:    translate.NewLobbySummaryTranslator(),
		LobbyPlayerSlot: translate.NewLobbyPlayerSlotSummaryTranslator(),
		LobbyPlayer:     translate.NewLobbyPlayerSummaryTranslator(),
	}
	application.subscribers = []ApplicationSubscriber{
		subscriber.NewLobbyNotificationChanneler(core.Services.EventPublisher, core.Services.Redis, core.Services.RabbitMQConnection, memory.NewLobbyChannelRepository(core.Services.Redis)),
		subscriber.NewLobbyPlayerSummaryWriter(core.Services.EventPublisher, core.data.LobbyPlayerSummary),
		subscriber.NewLobbySummaryWriter(core.Services.EventPublisher, core.repositories.Quests, core.data.LobbySummary),
		subscriber.NewLobbySearchWriter(core.Services.EventPublisher, core.data.LobbySearch),
		subscriber.NewLobbyChannelEventNotifier(core.Services.EventPublisher, core.data.LobbyPlayerSummary, application.Translators.LobbyPlayer),
		subscriber.NewLobbyClientNotifier(core.Services.EventPublisher, core.Services.Redis),
		subscriber.NewSessionLobbyWriter(core.Services.EventPublisher, core.data.SessionInstance),
	}
	return application
}

type CreateSessionCommandHandler decorator.CommandHandler[command.Context, *command.CreateSessionCommand]
type EndSessionCommandHandler decorator.CommandHandler[command.Context, *command.EndSessionCommand]
type CreateLobbyCommandHandler decorator.CommandHandler[command.Context, *command.CreateLobbyCommand]
type CancelLobbyCommandHandler decorator.CommandHandler[command.Context, *command.CancelLobbyCommand]
type WatchLobbyCommandHandler decorator.CommandHandler[command.Context, *command.WatchLobbyCommand]
type JoinLobbyCommandHandler decorator.CommandHandler[command.Context, *command.JoinLobbyCommand]
type LeaveLobbyCommandHandler decorator.CommandHandler[command.Context, *command.LeaveLobbyCommand]
type ReadyLobbyCommandHandler decorator.CommandHandler[command.Context, *command.ReadyLobbyCommand]
type UnreadyLobbyCommandHandler decorator.CommandHandler[command.Context, *command.UnreadyLobbyCommand]
type SendStampCommandHandler decorator.CommandHandler[command.Context, *command.SendStampCommand]

type SearchLobbyQueryHandler decorator.QueryHandler[query.Context, query.SearchLobbyQuery, []lobby.Summary]
type SearchPlayerQueryHandler decorator.QueryHandler[query.Context, query.SearchPlayerQuery, lobby.PlayerSummary]

func (a *LobbyApplication) NewSearchLobbyQueryHandler(core *CoreApplication) SearchLobbyQueryHandler {
	var actual = query.NewSearchLobbyQueryHandler(core.data.LobbySearch, lobby.NewSummaryQueryService(core.data.LobbyInstance, core.data.LobbyParticipant, core.data.LobbySummary, core.data.LobbyPlayerSummary))
	return actual
}

func (a *LobbyApplication) NewSearchPlayerQueryHandler(core *CoreApplication) SearchPlayerQueryHandler {
	var actual = query.NewSearchPlayerQueryHandler(core.data.LobbyParticipant, core.data.SessionInstance, core.data.LobbyPlayerSummary)
	return actual
}

func (a *LobbyApplication) NewCreateSessionCommandHandler(core *CoreApplication) CreateSessionCommandHandler {
	var actual = command.NewCreateSessionCommandHandler(core.Services.EventPublisher, core.data.SessionInstance)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CreateSessionCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewEndSessionCommandHandler(core *CoreApplication) EndSessionCommandHandler {
	var actual = command.NewEndSessionCommandHandler(core.Services.EventPublisher, core.data.SessionInstance, core.data.SessionInstance)
	return decorator.NewStandardCommandDecorator[command.Context, *command.EndSessionCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewCreateLobbyCommandHandler(core *CoreApplication) CreateLobbyCommandHandler {
	var actual = command.NewCreateLobbyCommandHandler(core.Services.EventPublisher, core.data.LobbyInstance, core.repositories.Quests, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CreateLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewCancelLobbyCommandHandler(core *CoreApplication) CancelLobbyCommandHandler {
	var actual = command.NewCancelLobbyCommandHandler(core.Services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CancelLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewWatchLobbyCommandHandler(core *CoreApplication) WatchLobbyCommandHandler {
	var actual = command.NewWatchLobbyCommandHandler(core.Services.EventPublisher)
	return decorator.NewStandardCommandDecorator[command.Context, *command.WatchLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewJoinLobbyCommandHandler(core *CoreApplication) JoinLobbyCommandHandler {
	var actual = command.NewJoinLobbyCommandHandler(core.Services.EventPublisher, core.data.LobbyInstance, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.JoinLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLeaveLobbyCommand(core *CoreApplication) LeaveLobbyCommandHandler {
	var actual = command.NewLeaveLobbyCommandHandler(core.Services.EventPublisher, core.data.SessionInstance, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LeaveLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewReadyLobbyCommandHandler(core *CoreApplication) ReadyLobbyCommandHandler {
	var actual = command.NewReadyLobbyCommandHandler(core.Services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ReadyLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewUnreadyLobbyCommandHandler(core *CoreApplication) UnreadyLobbyCommandHandler {
	var actual = command.NewUnreadyLobbyCommandHandler(core.Services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return decorator.NewStandardCommandDecorator[command.Context, *command.UnreadyLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewSendStampCommandHandler(core *CoreApplication) SendStampCommandHandler {
	var actual = command.NewSendStampCommandHandler(core.Services.EventPublisher, core.data.SessionInstance)
	return decorator.NewStandardCommandDecorator[command.Context, *command.SendStampCommand](core.Services.EventPublisher, actual)
}
