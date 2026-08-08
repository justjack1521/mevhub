package subscriber

import (
	"context"
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
		// Fully retire the channel — closing without deregistering leaked the
		// registry entry (and the subscription) for the process lifetime.
		c.manager.remove(ctx, c.LobbyID)
		return nil
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
			c.manager.logger.With("lobby.id", c.LobbyID.String(), "error", err.Error()).Error("failed to publish lobby notification to listener")
		}
	}
}

func (c *LobbyInstanceNotificationChannel) close(ctx context.Context) {
	// Close (not just Unsubscribe, which was previously issued with the wrong
	// channel name and did nothing) tears down the subscription and ends the
	// run() goroutine by closing its message channel.
	if err := c.channel.Close(); err != nil {
		c.manager.logger.With("lobby.id", c.LobbyID.String(), "error", err.Error()).Error("failed to close lobby notification channel")
	}
}
