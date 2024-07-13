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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"mevhub/internal/adapter/broker"
	"mevhub/internal/adapter/database"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/app"
	"mevhub/internal/ports"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	var logger = logrus.New()

	var application = NewApplication(context.Background(), logger)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic: %v", r)
			os.Exit(1)
			return
		}
		if err := application.Shutdown(); err != nil {
			fmt.Println(err)
		}
	}()

	var interceptions []grpc.UnaryServerInterceptor
	if application.Services.NewRelic != nil && application.Services.NewRelic.Application != nil {
		interceptions = append(interceptions, nrgrpc.UnaryServerInterceptor(application.Services.NewRelic.Application))
	}
	interceptions = append(interceptions, mevrpc.IdentityExtractionUnaryServerInterceptor)

	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptions...),
	}

	application.Start()

	server.RunGRPCServerWithOptions("50555", func(svr *grpc.Server) {
		svc := ports.NewMultiGrpcServer(application)
		services.RegisterMeviusMultiServiceServer(svr, svc)
	}, options...)

}

func NewApplication(ctx context.Context, logger *logrus.Logger) *app.CoreApplication {

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	rds, err := memory.NewRedisConnection(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to cache: %w", err))
	}

	msg, err := broker.NewRabbitMQConnection()
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

	return app.NewApplication(db, rds, logger, msg, identity, app.ApplicationWithNewRelic(nrl))

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
