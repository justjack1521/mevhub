package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrFailedConnectToPostgres = func(err error) error {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	ErrFailedConnectToRedis = func(err error) error {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
)

func NewPostgresConnection() (*gorm.DB, error) {
	config, err := mevconn.NewPostgresConfig()
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}
	db, err := gorm.Open(postgres.Open(config.Source()), &gorm.Config{})
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}
	return db, nil
}
func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	config, err := mevconn.NewRedisConfig()
	if err != nil {
		return nil, ErrFailedConnectToRedis(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.DSN(),
		Password: config.Password(),
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, ErrFailedConnectToRedis(err)
	}
	return client, nil
}
