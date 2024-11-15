package server

import (
	"context"
	"github.com/justjack1521/mevrabbit"
	"github.com/wagslane/go-rabbitmq"
)

type GameServerRabbitMQNotifier struct {
	publisher *mevrabbit.StandardPublisher
}

func NewGameServerRabbitMQNotifier(conn *rabbitmq.Conn) *GameServerRabbitMQNotifier {
	return &GameServerRabbitMQNotifier{
		publisher: mevrabbit.NewClientPublisher(conn),
	}
}

func (n *GameServerRabbitMQNotifier) Publish(ctx context.Context, player *PlayerChannel, notification Notification) error {
	data, err := notification.MarshallBinary()
	if err != nil {
		return err
	}
	return n.publisher.Publish(ctx, data, player.UserID, player.PlayerID, mevrabbit.ClientNotification)
}
