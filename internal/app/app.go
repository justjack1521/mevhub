package app

import (
	"github.com/go-redis/redis/v8"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/justjack1521/mevrelic"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
	"mevhub/internal/adapter/cache"
	"mevhub/internal/adapter/database"
	"mevhub/internal/adapter/external"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/adapter/serial"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
	"mevhub/internal/domain/session"
)

type CoreApplication struct {
	SubApplications *SubApplications
	repositories    *Repositories
	data            *DataRepositories
	Services        ApplicationServices
}

func (a *CoreApplication) Start() {
	a.Services.EventPublisher.Notify(mevent.ApplicationStartEvent{})
}

func (a *CoreApplication) Shutdown() error {
	a.Services.EventPublisher.Notify(mevent.ApplicationShutdownEvent{})
	return nil
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

type ApplicationServices struct {
	Logger             *logrus.Logger
	EventPublisher     *mevent.Publisher
	Redis              *redis.Client
	RabbitMQConnection *rabbitmq.Conn
	IdentityService    services.MeviusIdentityServiceClient
	NewRelic           *mevrelic.NewRelic
}

type CoreApplicationConfigurationOption func(c *CoreApplication) *CoreApplication

func New() *CoreApplication {
	return &CoreApplication{}
}

func NewApplication(db *gorm.DB, client *redis.Client, logger *logrus.Logger, conn *rabbitmq.Conn, identity services.MeviusIdentityServiceClient, options ...CoreApplicationConfigurationOption) *CoreApplication {
	var application = New().BuildServices(client, conn, logger, identity).BuildRepos(db, client).BuildDataRepos(db, client, identity).BuildSubApps()
	for _, opt := range options {
		opt(application)
	}
	return application
}

func (a *CoreApplication) BuildServices(client *redis.Client, mq *rabbitmq.Conn, logger *logrus.Logger, identity services.MeviusIdentityServiceClient) *CoreApplication {
	publisher := mevent.NewPublisher(mevent.PublisherWithLogger(logger))
	a.Services = ApplicationServices{
		Logger:             logger,
		EventPublisher:     publisher,
		Redis:              client,
		RabbitMQConnection: mq,
		IdentityService:    identity,
	}
	return a
}

func (a *CoreApplication) BuildRepos(db *gorm.DB, client *redis.Client) *CoreApplication {
	a.repositories = &Repositories{
		Quests: database.NewGameQuestDatabaseRepository(db),
	}
	return a
}

func (a *CoreApplication) BuildDataRepos(db *gorm.DB, client *redis.Client, identity services.MeviusIdentityServiceClient) *CoreApplication {
	a.data = &DataRepositories{
		SessionInstance:    memory.NewLobbySessionRedisRepository(client),
		LobbyInstance:      memory.NewLobbyInstanceRedisRepository(client),
		LobbyParticipant:   memory.NewLobbyParticipantRedisRepository(client),
		LobbySearch:        memory.NewLobbySearchRepository(client),
		LobbySummary:       database.NewLobbySummaryDatabaseRepository(db),
		LobbyPlayerSummary: cache.NewLobbyPlayerSummaryRepository(external.NewLobbyPlayerSummaryRepository(identity), memory.NewPlayerSummaryRepository(client, serial.LobbyPlayerSummaryJSONSerialiser{})),
		GameInstance:       database.NewGameInstanceDatabaseRepository(db),
		GameSummary:        nil,
		GamePlayerSummary:  nil,
	}
	return a
}

func (a *CoreApplication) BuildSubApps() *CoreApplication {
	a.SubApplications = &SubApplications{}
	a.SubApplications.Game = NewGameApplication(a)
	a.SubApplications.Lobby = NewLobbyApplication(a)
	return a
}
