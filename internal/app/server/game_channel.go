package server

import (
	"context"
	"fmt"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/rabbitmv"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/domain/game"
	"sync"
)

type GameChannel struct {
	mu          sync.Mutex
	InstanceID  uuid.UUID
	instance    *game.Instance
	players     map[uuid.UUID]*PlayerChannel
	Register    chan *PlayerChannel
	Unregister  chan uuid.UUID
	ReadyPlayer chan *PlayerReadyNotification

	Publisher *rabbitmv.StandardPublisher
}

func NewGameChannel(instance *game.Instance, conn *rabbitmq.Conn, logger *logrus.Logger) *GameChannel {
	return &GameChannel{
		InstanceID:  instance.SysID,
		instance:    instance,
		players:     make(map[uuid.UUID]*PlayerChannel),
		Register:    make(chan *PlayerChannel),
		Unregister:  make(chan uuid.UUID),
		ReadyPlayer: make(chan *PlayerReadyNotification),
		Publisher:   rabbitmv.NewClientPublisher(conn, rabbitmq.WithPublisherOptionsLogger(logger)),
	}
}

func (c *GameChannel) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.Register)
	close(c.Unregister)
	close(c.ReadyPlayer)
}

func (c *GameChannel) Run() {
	for {
		select {
		case player := <-c.Register:
			c.register(player)
		case id := <-c.Unregister:
			c.unregister(id)
		case notification := <-c.ReadyPlayer:
			c.readyPlayer(notification)
		}
	}
}

func (c *GameChannel) register(channel *PlayerChannel) {
	c.players[channel.PlayerID] = channel
}

func (c *GameChannel) unregister(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.players[id]; exists {
		delete(c.players, id)
	}
}

func (c *GameChannel) readyPlayer(notification *PlayerReadyNotification) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if instance, exists := c.players[notification.InstanceID]; exists {
		instance.Ready = true
		var message = &protomulti.GamePlayerReadyNotification{
			GameId:      c.InstanceID.String(),
			PlayerIndex: int32(instance.participant.PlayerSlot),
		}
		if err := c.publish(context.Background(), protomulti.MultiGameNotificationType_GAME_PLAYER_READY, message); err != nil {
			fmt.Println(err)
		}
	}

	for _, player := range c.players {
		if player.Ready == false {
			break
		}
		c.instance.State = game.InstanceGameStartedState
		var message = &protomulti.GameReadyNotification{GameId: c.InstanceID.String()}
		if err := c.publish(context.Background(), protomulti.MultiGameNotificationType_GAME_READY, message); err != nil {
			fmt.Println(err)
		}
	}

}

func (c *GameChannel) publish(ctx context.Context, operation protomulti.MultiGameNotificationType, message Notification) error {

	bytes, err := message.MarshallBinary()
	if err != nil {
		return err
	}

	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(operation),
		Data:    bytes,
	}

	msg, err := notification.MarshallBinary()
	if err != nil {
		return err
	}

	for _, player := range c.players {
		if err := c.Publisher.Publish(ctx, msg, player.UserID, player.PlayerID, rabbitmv.ClientNotification); err != nil {
			return err
		}
	}
	return nil
}

type Notification interface {
	MarshallBinary() ([]byte, error)
}
