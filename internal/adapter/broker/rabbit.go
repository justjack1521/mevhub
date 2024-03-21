package broker

import (
	"fmt"
	"github.com/justjack1521/mevconn"
	"github.com/wagslane/go-rabbitmq"
	"os"
)

var (
	ErrEnvironmentVariableMissing = func(env string) error {
		return fmt.Errorf("environment variable %s is missing", env)
	}
	ErrFailedConnectToRabbitMQ = func(err error) error {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
)

func buildURL() (string, error) {
	var username = os.Getenv("RMQUSERNAME")
	if username == "" {
		return "", ErrEnvironmentVariableMissing("RMQUSERNAME")
	}
	var password = os.Getenv("RMQPASSWORD")
	if username == "" {
		return "", ErrEnvironmentVariableMissing("RMQPASSWORD")
	}
	var host = os.Getenv("RMQHOST")
	if username == "" {
		return "", ErrEnvironmentVariableMissing("RMQHOST")
	}
	var port = os.Getenv("RMQPORT")
	if username == "" {
		return "", ErrEnvironmentVariableMissing("RMQPORT")
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port), nil
}

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
