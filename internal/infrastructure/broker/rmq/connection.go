package rmq

import (
	"fmt"
	"github.com/justjack1521/mevconn"
	"github.com/wagslane/go-rabbitmq"
)

var (
	ErrFailedConnectToRabbitMQ = func(err error) error {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
)

func NewRabbitMQConnection() (*rabbitmq.Conn, error) {
	config, err := mevconn.CreateRabbitMQConfig()
	if err != nil {
		return nil, ErrFailedConnectToRabbitMQ(err)
	}
	conn, err := rabbitmq.NewConn(config.Source(), rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, ErrFailedConnectToRabbitMQ(err)
	}
	return conn, nil
}
