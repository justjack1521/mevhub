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
	services    *Services
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

func NewLobbyApplication(core *Application) *LobbyApplication {
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
		subscriber.NewLobbyNotificationChanneler(core.services.EventPublisher, core.services.Redis, core.services.RabbitMQConnection, memory.NewLobbyChannelRepository(core.services.Redis)),
		subscriber.NewLobbyPlayerSummaryWriter(core.services.EventPublisher, core.data.LobbyPlayerSummary),
		subscriber.NewLobbySummaryWriter(core.services.EventPublisher, core.repositories.Quests, core.data.LobbySummary),
		subscriber.NewLobbySearchWriter(core.services.EventPublisher, core.data.LobbySearch),
		subscriber.NewLobbyChannelEventNotifier(core.services.EventPublisher, core.data.LobbyPlayerSummary, application.Translators.LobbyPlayer),
		subscriber.NewLobbyClientNotifier(core.services.EventPublisher, core.services.Redis),
		subscriber.NewSessionLobbyWriter(core.services.EventPublisher, core.data.SessionInstance),
	}
	return application
}

type SearchLobbyQueryHandler decorator.QueryHandler[query.SearchLobbyQuery, []lobby.Summary]

func (a *LobbyApplication) NewSearchLobbyQueryHandler(core *Application) SearchLobbyQueryHandler {
	var actual = query.NewSearchLobbyQueryHandler(core.data.LobbySearch, lobby.NewSummaryQueryService(core.data.LobbyInstance, core.data.LobbyParticipant, core.data.LobbySummary, core.data.LobbyPlayerSummary))
	return ApplyStandardQueryDecorators[query.SearchLobbyQuery, []lobby.Summary](core, actual)
}

type SearchPlayerQueryHandler decorator.QueryHandler[query.SearchPlayerQuery, lobby.PlayerSummary]

func (a *LobbyApplication) NewSearchPlayerQueryHandler(core *Application) SearchPlayerQueryHandler {
	var actual = query.NewSearchPlayerQueryHandler(core.data.LobbyParticipant, core.data.SessionInstance, core.data.LobbyPlayerSummary)
	return ApplyStandardQueryDecorators[query.SearchPlayerQuery, lobby.PlayerSummary](core, actual)
}

type CreateSessionCommandHandler decorator.CommandHandler[command.CreateSessionCommand]

func (a *LobbyApplication) NewCreateSessionCommandHandler(core *Application) CreateSessionCommandHandler {
	var actual = command.NewCreateSessionCommandHandler(core.services.EventPublisher, core.data.SessionInstance)
	var logger = decorator.NewCommandHandlerWithLogger[command.CreateSessionCommand](core.services.Logger, actual)
	return logger
}

type EndSessionCommandHandler decorator.CommandHandler[command.EndSessionCommand]

func (a *LobbyApplication) NewEndSessionCommandHandler(core *Application) EndSessionCommandHandler {
	var actual = command.NewEndSessionCommandHandler(core.services.EventPublisher, core.data.SessionInstance, core.data.SessionInstance)
	var logger = decorator.NewCommandHandlerWithLogger[command.EndSessionCommand](core.services.Logger, actual)
	return logger
}

type CreateLobbyCommandHandler decorator.CommandHandler[command.CreateLobbyCommand]

func (a *LobbyApplication) NewCreateLobbyCommandHandler(core *Application) CreateLobbyCommandHandler {
	var actual = command.NewCreateLobbyCommandHandler(core.services.EventPublisher, core.data.LobbyInstance, core.repositories.Quests, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.CreateLobbyCommand](core, actual)
}

type CancelLobbyCommandHandler decorator.CommandHandler[command.CancelLobbyCommand]

func (a *LobbyApplication) NewCancelLobbyCommandHandler(core *Application) CancelLobbyCommandHandler {
	var actual = command.NewCancelLobbyCommandHandler(core.services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.CancelLobbyCommand](core, actual)
}

type WatchLobbyCommandHandler decorator.CommandHandler[command.WatchLobbyCommand]

func (a *LobbyApplication) NewWatchLobbyCommandHandler(core *Application) WatchLobbyCommandHandler {
	var actual = command.NewWatchLobbyCommandHandler(core.services.EventPublisher)
	return ApplyStandardCommandDecorators[command.WatchLobbyCommand](core, actual)
}

type JoinLobbyCommandHandler decorator.CommandHandler[command.JoinLobbyCommand]

func (a *LobbyApplication) NewJoinLobbyCommandHandler(core *Application) JoinLobbyCommandHandler {
	var actual = command.NewJoinLobbyCommandHandler(core.services.EventPublisher, core.data.LobbyInstance, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.JoinLobbyCommand](core, actual)
}

type LeaveLobbyCommandHandler decorator.CommandHandler[command.LeaveLobbyCommand]

func (a *LobbyApplication) NewLeaveLobbyCommand(core *Application) LeaveLobbyCommandHandler {
	var actual = command.NewLeaveLobbyCommandHandler(core.services.EventPublisher, core.data.SessionInstance, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.LeaveLobbyCommand](core, actual)
}

type ReadyLobbyCommandHandler decorator.CommandHandler[command.ReadyLobbyCommand]

func (a *LobbyApplication) NewReadyLobbyCommandHandler(core *Application) ReadyLobbyCommandHandler {
	var actual = command.NewReadyLobbyCommandHandler(core.services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.ReadyLobbyCommand](core, actual)
}

type UnreadyLobbyCommandHandler decorator.CommandHandler[command.UnreadyLobbyCommand]

func (a *LobbyApplication) NewUnreadyLobbyCommandHandler(core *Application) UnreadyLobbyCommandHandler {
	var actual = command.NewUnreadyLobbyCommandHandler(core.services.EventPublisher, core.data.SessionInstance, core.data.LobbyInstance, core.data.LobbyParticipant)
	return ApplyStandardCommandDecorators[command.UnreadyLobbyCommand](core, actual)
}

type SendStampCommandHandler decorator.CommandHandler[command.SendStampCommand]

func (a *LobbyApplication) NewSendStampCommandHandler(core *Application) SendStampCommandHandler {
	var actual = command.NewSendStampCommandHandler(core.services.EventPublisher, core.data.SessionInstance)
	return ApplyStandardCommandDecorators[command.SendStampCommand](core, actual)
}
