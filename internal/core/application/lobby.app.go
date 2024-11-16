package application

import (
	"mevhub/internal/adapter/memory"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/application/command"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/application/query"
	"mevhub/internal/core/application/subscriber"
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
	CreateSession CreateSessionCommandHandler
	EndSession    EndSessionCommandHandler
	CreateLobby   CreateLobbyCommandHandler
	WatchLobby    WatchLobbyCommandHandler
	JoinLobby     JoinLobbyCommandHandler
	CancelLobby   CancelLobbyCommandHandler
	ReadyLobby    ReadyParticipantCommandHandler
	UnreadyLobby  UnreadyParticipantCommandHandler
	StartLobby    StartLobbyCommandHandler
	SendStamp     SendStampCommandHandler
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
		CreateSession: application.NewCreateSessionCommandHandler(core),
		EndSession:    application.NewEndSessionCommandHandler(core),
		CreateLobby:   application.NewCreateLobbyCommandHandler(core),
		CancelLobby:   application.NewCancelLobbyCommandHandler(core),
		WatchLobby:    application.NewWatchLobbyCommandHandler(core),
		JoinLobby:     application.NewJoinLobbyCommandHandler(core),
		ReadyLobby:    application.NewReadyParticipantCommandHandler(core),
		UnreadyLobby:  application.NewUnreadyParticipantCommandHandler(core),
		StartLobby:    application.NewStartLobbyCommandHandler(core),
		SendStamp:     application.NewSendStampCommandHandler(core),
	}
	application.Translators = &LobbyApplicationTranslators{
		LobbySummary:    translate.NewLobbySummaryTranslator(),
		LobbyPlayerSlot: translate.NewLobbyPlayerSlotSummaryTranslator(),
		LobbyPlayer:     translate.NewLobbyPlayerSummaryTranslator(),
	}
	application.subscribers = []ApplicationSubscriber{
		subscriber.NewLobbyNotificationChanneler(core.Services.EventPublisher, core.Services.Redis, core.Services.RabbitMQConnection, memory.NewLobbyChannelRepository(core.Services.Redis)),
		subscriber.NewLobbySummaryWriter(core.Services.EventPublisher, core.repositories.Quests, core.data.LobbySummaries),
		subscriber.NewLobbySearchWriter(core.Services.EventPublisher, core.data.LobbySearch),
		subscriber.NewLobbyChannelEventNotifier(core.Services.EventPublisher, core.data.LobbyPlayerSummaries, application.Translators.LobbyPlayer),
		subscriber.NewLobbyClientNotifier(core.Services.EventPublisher, core.Services.Redis),
		subscriber.NewSessionLobbyWriter(core.Services.EventPublisher, core.data.Sessions),
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
type ReadyParticipantCommandHandler decorator.CommandHandler[command.Context, *command.ReadyParticipantCommand]
type UnreadyParticipantCommandHandler decorator.CommandHandler[command.Context, *command.UnreadyParticipantCommand]
type StartLobbyCommandHandler decorator.CommandHandler[command.Context, *command.StartLobbyCommand]
type SendStampCommandHandler decorator.CommandHandler[command.Context, *command.SendStampCommand]

type SearchLobbyQueryHandler decorator.QueryHandler[query.Context, query.SearchLobbyQuery, []lobby.Summary]
type SearchPlayerQueryHandler decorator.QueryHandler[query.Context, query.SearchPlayerQuery, lobby.PlayerSummary]

func (a *LobbyApplication) NewSearchLobbyQueryHandler(core *CoreApplication) SearchLobbyQueryHandler {
	var actual = query.NewSearchLobbyQueryHandler(core.data.LobbySearch, lobby.NewSummaryQueryService(core.data.Lobbies, core.data.LobbyParticipants, core.data.LobbySummaries, core.data.LobbyPlayerSummaries))
	return actual
}

func (a *LobbyApplication) NewSearchPlayerQueryHandler(core *CoreApplication) SearchPlayerQueryHandler {
	var actual = query.NewSearchPlayerQueryHandler(core.data.LobbyParticipants, core.data.Sessions, core.data.LobbyPlayerSummaries)
	return actual
}

func (a *LobbyApplication) NewCreateSessionCommandHandler(core *CoreApplication) CreateSessionCommandHandler {
	var actual = command.NewCreateSessionCommandHandler(core.Services.EventPublisher, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CreateSessionCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewEndSessionCommandHandler(core *CoreApplication) EndSessionCommandHandler {
	var actual = command.NewEndSessionCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.EndSessionCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewCreateLobbyCommandHandler(core *CoreApplication) CreateLobbyCommandHandler {
	var actual = command.NewCreateLobbyCommandHandler(core.Services.EventPublisher, core.data.Lobbies, core.repositories.Quests, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CreateLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewCancelLobbyCommandHandler(core *CoreApplication) CancelLobbyCommandHandler {
	var actual = command.NewCancelLobbyCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.CancelLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewStartLobbyCommandHandler(core *CoreApplication) StartLobbyCommandHandler {
	var actual = command.NewStartLobbyCommandHandler(core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.StartLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewWatchLobbyCommandHandler(core *CoreApplication) WatchLobbyCommandHandler {
	var actual = command.NewWatchLobbyCommandHandler(core.Services.EventPublisher)
	return decorator.NewStandardCommandDecorator[command.Context, *command.WatchLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewJoinLobbyCommandHandler(core *CoreApplication) JoinLobbyCommandHandler {
	var actual = command.NewJoinLobbyCommandHandler(core.Services.EventPublisher, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.JoinLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewLeaveLobbyCommand(core *CoreApplication) LeaveLobbyCommandHandler {
	var actual = command.NewLeaveLobbyCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.LeaveLobbyCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewReadyParticipantCommandHandler(core *CoreApplication) ReadyParticipantCommandHandler {
	var actual = command.NewReadyParticipantCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.ReadyParticipantCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewUnreadyParticipantCommandHandler(core *CoreApplication) UnreadyParticipantCommandHandler {
	var actual = command.NewUnreadyParticipantCommandHandler(core.Services.EventPublisher, core.data.Sessions, core.data.Lobbies, core.data.LobbyParticipants)
	return decorator.NewStandardCommandDecorator[command.Context, *command.UnreadyParticipantCommand](core.Services.EventPublisher, actual)
}

func (a *LobbyApplication) NewSendStampCommandHandler(core *CoreApplication) SendStampCommandHandler {
	var actual = command.NewSendStampCommandHandler(core.Services.EventPublisher, core.data.Sessions)
	return decorator.NewStandardCommandDecorator[command.Context, *command.SendStampCommand](core.Services.EventPublisher, actual)
}
