package config

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevcon"
)

type Application struct {
	Database   mevcon.PostgresConfig `required:"true"`
	Redis      mevcon.RedisConfig    `required:"true"`
	RabbitMQ   mevcon.RabbitMQConfig `required:"true"`
	GameClient Client                `required:"true"`
}

type Client struct {
	Address         string `required:"true"`
	Port            string `required:"true"`
	CertificatePath string
}

func (c Client) ConnectionString() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}
