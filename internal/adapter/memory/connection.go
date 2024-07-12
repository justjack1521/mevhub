package memory

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevconn"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
)

const (
	serviceKey = "mevmulti"
)

var (
	ErrFailedConnectToRedis = func(err error) error {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
)

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	config, err := mevconn.NewRedisConfig()
	if err != nil {
		return nil, ErrFailedConnectToRedis(err)
	}
	var options = &redis.Options{
		Addr:      config.DSN(),
		Username:  config.Username(),
		Password:  config.Password(),
		TLSConfig: &tls.Config{ServerName: config.Host()},
	}
	client := redis.NewClient(options)
	client.AddHook(nrredis.NewHook(options))
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, ErrFailedConnectToRedis(err)
	}
	return client, nil
}
