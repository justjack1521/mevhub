package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection() (*gorm.DB, error) {
	config, err := mevconn.NewPostgresConfig()
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(config.Source()), &gorm.Config{})
}

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	config, err := mevconn.NewRedisConfig()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.DSN(),
		Password: config.Password(),
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
