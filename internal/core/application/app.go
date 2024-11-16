package application

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
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
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
	Options port.InstanceOptionsRepository
	Quests  port.QuestRepository
}

type DataRepositories struct {
	Sessions             session.InstanceRepository
	Lobbies              port.LobbyInstanceRepository
	LobbyParticipants    lobby.ParticipantRepository
	LobbySummaries       lobby.SummaryRepository
	LobbyPlayerSummaries lobby.PlayerSummaryRepository
	LobbySearch          lobby.SearchRepository
	Games                port.GameInstanceRepository
	GameParticipants     port.PlayerParticipantRepository
	GamePlayerLoadouts   port.PlayerLoadoutReadRepository
}

type ApplicationServices struct {
	Logger             *logrus.Logger
	EventPublisher     *mevent.Publisher
	Redis              *redis.Client
	RabbitMQConnection *rabbitmq.Conn
	IdentityService    services.MeviusIdentityServiceClient
	NewRelic           *mevrelic.NewRelic
}

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
		Sessions:             memory.NewLobbySessionRedisRepository(client),
		Lobbies:              memory.NewLobbyInstanceRedisRepository(client),
		LobbyParticipants:    memory.NewLobbyParticipantRedisRepository(client),
		LobbySearch:          memory.NewLobbySearchRepository(client),
		LobbySummaries:       database.NewLobbySummaryDatabaseRepository(db),
		LobbyPlayerSummaries: cache.NewLobbyPlayerSummaryRepository(external.NewLobbyPlayerSummaryRepository(identity), memory.NewLobbyPlayerSummaryRepository(client, serial.NewLobbyPlayerSummaryJSONSerialiser())),
		Games:                memory.NewGameInstanceRepository(client, serial.NewGameInstanceJSONSerialiser()),
		GameParticipants:     memory.NewGameParticipantRepository(client, serial.NewGamePlayerParticipantJSONSerialiser()),
		GamePlayerLoadouts:   external.NewGamePlayerLoadoutRepository(identity),
	}
	return a
}

func (a *CoreApplication) BuildSubApps() *CoreApplication {
	a.SubApplications = &SubApplications{}
	a.SubApplications.Game = NewGameApplication(a)
	a.SubApplications.Lobby = NewLobbyApplication(a)
	return a
}

type CoreApplicationConfigurationOption func(c *CoreApplication) *CoreApplication

func ApplicationWithNewRelic(relic *mevrelic.NewRelic) CoreApplicationConfigurationOption {
	return func(c *CoreApplication) *CoreApplication {
		c.Services.NewRelic = relic
		if c.Services.Logger != nil {
			c.Services.NewRelic.Attach(c.Services.Logger)
		}
		return c
	}
}
