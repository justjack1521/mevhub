package app

import (
	"github.com/go-redis/redis/v8"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
	"mevhub/internal/adapter/cache"
	"mevhub/internal/adapter/database"
	"mevhub/internal/adapter/extern"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/decorator"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type Application struct {
	SubApplications *SubApplications
	repositories    *Repositories
	data            *DataRepositories
	services        *Services
}

type SubApplications struct {
	Game  *GameApplication
	Lobby *LobbyApplication
}

type Repositories struct {
	Options game.InstanceOptionsRepository
	Quests  game.QuestRepository
}

type DataRepositories struct {
	SessionInstance    session.InstanceRepository
	GameInstance       game.InstanceRepository
	GameSummary        game.SummaryRepository
	LobbyParticipant   lobby.ParticipantRepository
	GamePlayerSummary  game.PlayerSummaryRepository
	LobbyInstance      lobby.InstanceRepository
	LobbySearch        lobby.SearchRepository
	LobbySummary       lobby.SummaryRepository
	LobbyPlayerSummary lobby.PlayerSummaryRepository
}

type Services struct {
	Logger             *logrus.Logger
	EventPublisher     *mevent.Publisher
	Redis              *redis.Client
	RabbitMQConnection *rabbitmq.Conn
}

func NewApplication(db *gorm.DB, client *redis.Client, logger *logrus.Logger, conn *rabbitmq.Conn, game services.MeviusGameServiceClient) *Application {
	var application = &Application{
		repositories: &Repositories{
			Quests: database.NewGameQuestDatabaseRepository(db),
		},
		data: &DataRepositories{
			SessionInstance:    memory.NewLobbySessionRedisRepository(client),
			LobbyInstance:      memory.NewLobbyInstanceRedisRepository(client),
			LobbyParticipant:   memory.NewLobbyParticipantRedisRepository(client),
			LobbySearch:        memory.NewLobbySearchRepository(client),
			LobbySummary:       database.NewLobbySummaryDatabaseRepository(db),
			LobbyPlayerSummary: cache.NewLazyLoadedLobbyPlayerSummaryRepository(extern.NewLobbyPlayerSummaryExternalRepository(game), memory.NewLobbyPlayerSlotSummaryRepository(client)),
			GameInstance:       database.NewGameInstanceDatabaseRepository(db),
			GameSummary:        nil,
			GamePlayerSummary:  extern.NewGamePlayerSummaryExternalRepository(game),
		},
		services: &Services{
			Logger:             logger,
			EventPublisher:     mevent.NewPublisher(mevent.PublisherWithLogger(logger)),
			Redis:              client,
			RabbitMQConnection: conn,
		},
	}
	application.SubApplications = &SubApplications{}
	application.SubApplications.Game = NewGameApplication(application)
	application.SubApplications.Lobby = NewLobbyApplication(application)
	return application
}

func (a *Application) Start() {
	a.services.EventPublisher.Notify(mevent.ApplicationStartEvent{})
}

func (a *Application) Shutdown() error {
	a.services.EventPublisher.Notify(mevent.ApplicationShutdownEvent{})
	return nil
}

func ApplyStandardCommandDecorators[C decorator.Command](core *Application, actual decorator.CommandHandler[C]) decorator.CommandHandler[C] {
	var logger = decorator.NewCommandHandlerWithLogger[C](core.services.Logger, actual)
	var qualifier = NewCommandHandlerWithSession[C](core.data.SessionInstance, logger)
	return qualifier
}

func ApplyStandardQueryDecorators[Q decorator.Query, R any](core *Application, actual decorator.QueryHandler[Q, R]) decorator.QueryHandler[Q, R] {
	var logger = decorator.NewQueryHandlerWithLogger[Q, R](core.services.Logger, actual)
	var qualifier = NewQueryHandlerWithSession[Q, R](core.data.SessionInstance, logger)
	return qualifier
}
