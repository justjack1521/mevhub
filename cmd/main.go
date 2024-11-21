package main

import (
	"context"
	"fmt"
	"github.com/justjack1521/mevconn"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	"github.com/justjack1521/mevium/pkg/server"
	"github.com/justjack1521/mevrelic"
	"github.com/justjack1521/mevrpc"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	slogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log/slog"
	"math/rand"
	"mevhub/internal/adapter/database"
	"mevhub/internal/adapter/handler/rpc"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/core/application"
	"mevhub/internal/infrastructure/broker/rmq"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	var logger = logrus.New()
	var slogger = slog.New(slogrus.Option{Level: slog.LevelDebug, Logger: logger}.NewLogrusHandler())

	var app = NewApplication(context.Background(), slogger)
	if app.Services.NewRelic != nil {
		app.Services.NewRelic.Attach(logger)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic: %v", r)
			os.Exit(1)
			return
		}
		if err := app.Shutdown(); err != nil {
			fmt.Println(err)
		}
	}()

	options := []grpc.ServerOption{
		//grpc.Creds(configuration.Certificates.NewTransportCredentials()),
		grpc.ChainUnaryInterceptor(
			nrgrpc.UnaryServerInterceptor(app.Services.NewRelic.Application),
			mevrpc.IdentityExtractionUnaryServerInterceptor,
		),
	}

	app.Start()

	server.RunGRPCServerWithOptions("50555", func(svr *grpc.Server) {
		svc := rpc.NewMultiGrpcServer(app)
		services.RegisterMeviusMultiServiceServer(svr, svc)
	}, options...)

}

func NewApplication(ctx context.Context, logger *slog.Logger) *application.CoreApplication {

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	rds, err := memory.NewRedisConnection(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to cache: %w", err))
	}

	msg, err := rmq.NewRabbitMQConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to message bus: %w", err))
	}

	identity, err := DialToIdentityClient()
	if err != nil {
		panic(fmt.Errorf("failed to connect to game client: %w", err))
	}

	nrl, err := mevrelic.NewRelicApplication()
	if err != nil {
		panic(err)
	}

	return application.NewApplication(db, rds, logger, msg, identity, application.ApplicationWithNewRelic(nrl))

}

func DialToIdentityClient() (services.MeviusIdentityServiceClient, error) {
	config, err := mevconn.CreateGrpcServiceConfig(mevconn.GAMESERVICENAME)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(config.ConnectionString(), grpc.WithChainUnaryInterceptor(nrgrpc.UnaryClientInterceptor, mevrpc.IdentityCopyUnaryClientInterceptor), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("Connected to %s", config.ConnectionString()))
	return services.NewMeviusIdentityServiceClient(conn), nil
}
