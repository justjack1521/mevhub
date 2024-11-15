package subscriber

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/justjack1521/mevrabbit"
	uuid "github.com/satori/go.uuid"
)

type LobbyInstanceNotificationChannel struct {
	LobbyID uuid.UUID
	channel *redis.PubSub
	manager *LobbyNotificationChanneler
}

func (c *LobbyInstanceNotificationChannel) Publish(ctx context.Context, message []byte) error {

	listeners, err := c.manager.repository.QueryAllForLobby(ctx, c.LobbyID)

	if err != nil {
		return err
	}

	if len(listeners) == 0 {
		c.close(ctx)
	}

	for _, listener := range listeners {
		if err := c.manager.publisher.Publish(ctx, message, listener.UserID, listener.PlayerID, mevrabbit.ClientNotification); err != nil {
			return err
		}
	}
	return nil
}

func (c *LobbyInstanceNotificationChannel) run() {
	channel := c.channel.Channel()
	for message := range channel {
		if err := c.Publish(context.Background(), []byte(message.Payload)); err != nil {
			fmt.Println("Error Publishing Message To Listener: ", err)
		}
	}
}

func (c *LobbyInstanceNotificationChannel) close(ctx context.Context) {
	if err := c.channel.Unsubscribe(ctx, c.LobbyID.String()); err != nil {
		fmt.Println(err)
	}
}
