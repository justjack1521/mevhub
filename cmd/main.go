package main

import (
	"context"
	"fmt"
	"github.com/justjack1521/mevconn"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	"github.com/justjack1521/mevium/pkg/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"mevhub/internal/adapter/broker"
	"mevhub/internal/adapter/database"
	"mevhub/internal/adapter/memory"
	"mevhub/internal/app"
	"mevhub/internal/ports"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	var logger = logrus.New()

	var application = NewApplication(context.Background(), logger)

	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(ports.ServerInterceptor(logger)),
	}

	application.Start()

	server.RunGRPCServerWithOptions("50555", func(svr *grpc.Server) {
		svc := ports.NewMultiGrpcServer(application)
		services.RegisterMeviusMultiServiceServer(svr, svc)
	}, options...)

	defer func(application *app.Application) {
		if err := application.Shutdown(); err != nil {
			fmt.Println(err)
		}
	}(application)

}

func NewApplication(ctx context.Context, logger *logrus.Logger) *app.Application {

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

	game, err := DialToGameClient()
	if err != nil {
		panic(fmt.Errorf("failed to connect to game client: %w", err))
	}

	return app.NewApplication(db, rds, logger, msg, game)

}

func DialToGameClient() (services.MeviusGameServiceClient, error) {
	config, err := mevconn.CreateGrpcServiceConfig(mevconn.GAMESERVICENAME)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(config.ConnectionString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return services.NewMeviusGameServiceClient(conn), nil
}
