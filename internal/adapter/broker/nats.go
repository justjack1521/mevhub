package broker

import (
	"github.com/justjack1521/mevconn"
	"github.com/nats-io/nats.go"
)

func NewNATSConnection() (*nats.Conn, error) {
	config, err := mevconn.NewNATSConfig()
	if err != nil {
		return nil, err
	}
	return nats.Connect(config.URL(), nats.Token(config.Token()))
}
