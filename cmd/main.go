package main

import (
	"context"
	"fmt"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	"github.com/justjack1521/mevium/pkg/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"mevhub/internal/adapter/broker"
	"mevhub/internal/adapter/database"
	"mevhub/internal/app"
	"mevhub/internal/config"
	"mevhub/internal/ports"
	"time"
)

const conf string = "/src/mevhub/internal/config/config.dev.json"

func main() {

	rand.Seed(time.Now().Unix())

	//var path = os.Getenv("GOPATH")
	var configuration config.Application
	//if err := configor.Load(&configuration, path+"/"+conf); err != nil {
	//	panic(err)
	//}

	var logger = logrus.New()

	var application = NewApplication(context.Background(), configuration, logger)

	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(ports.ServerInterceptor(logger)),
	}

	application.Start()

	server.RunGRPCServerWithOptions("50552", func(svr *grpc.Server) {
		svc := ports.NewMultiGrpcServer(application)
		services.RegisterMeviusMultiServiceServer(svr, svc)
	}, options...)

	defer func(application *app.Application) {
		if err := application.Shutdown(); err != nil {
			fmt.Println(err)
		}
	}(application)

}

func NewApplication(ctx context.Context, configuration config.Application, logger *logrus.Logger) *app.Application {

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	rds, err := database.NewRedisConnection(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to cache: %w", err))
	}

	msg, err := broker.NewRabbitMQConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to message bus: %w", err))
	}

	nts, err := broker.NewNATSConnection()
	if err != nil {
		panic(fmt.Errorf("failed to connect to nats: %w", err))
	}

	game, err := DialToGameClient(configuration)
	if err != nil {
		panic(fmt.Errorf("failed to connect to game client: %w", err))
	}

	return app.NewApplication(db, rds, logger, msg, nts, game)

}

func DialToGameClient(config config.Application) (services.MeviusGameServiceClient, error) {
	conn, err := grpc.Dial(config.GameClient.ConnectionString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return services.NewMeviusGameServiceClient(conn), nil
}
