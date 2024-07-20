package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/domain/game"
	"sync"
)

type GameChannelServer struct {
	mu         sync.Mutex
	games      map[uuid.UUID]*GameChannel
	Register   chan *GameChannel
	Unregister chan uuid.UUID

	connection *rabbitmq.Conn
	logger     *logrus.Logger

	RegisterPlayer chan *PlayerRegisterNotification
	ReadyPlayer    chan *PlayerReadyNotification
}

func NewGameChannelServer(conn *rabbitmq.Conn, logger *logrus.Logger) *GameChannelServer {
	var server = &GameChannelServer{
		connection: conn,
		logger:     logger,
		games:      make(map[uuid.UUID]*GameChannel),

		Register:       make(chan *GameChannel),
		Unregister:     make(chan uuid.UUID),
		RegisterPlayer: make(chan *PlayerRegisterNotification),
		ReadyPlayer:    make(chan *PlayerReadyNotification),
	}
	return server
}

func (c *GameChannelServer) Run() {

	for {
		select {
		case instance := <-c.Register:
			c.register(instance)
		case id := <-c.Unregister:
			c.unregister(id)
		case player := <-c.RegisterPlayer:
			c.registerPlayer(player)
		case player := <-c.ReadyPlayer:
			c.readyPlayer(player)
		}

	}

}

func (c *GameChannelServer) Close() {

	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.Register)
	close(c.Unregister)
	close(c.RegisterPlayer)
	close(c.ReadyPlayer)
	for _, instance := range c.games {
		instance.Close()
	}

}

func (c *GameChannelServer) NewGameChannel(instance *game.Instance) *GameChannel {
	return NewGameChannel(instance, c.connection, c.logger)
}

func (c *GameChannelServer) register(channel *GameChannel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.games[channel.InstanceID] = channel
	go channel.Run()
}

func (c *GameChannelServer) unregister(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if instance, exists := c.games[id]; exists {
		delete(c.games, id)
		instance.Close()
	}
}

func (c *GameChannelServer) registerPlayer(notification *PlayerRegisterNotification) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if instance, exists := c.games[notification.InstanceID]; exists {
		instance.Register <- notification.Player
	}
}

func (c *GameChannelServer) readyPlayer(notification *PlayerReadyNotification) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if instance, exists := c.games[notification.InstanceID]; exists {
		instance.ReadyPlayer <- notification
	}
}
