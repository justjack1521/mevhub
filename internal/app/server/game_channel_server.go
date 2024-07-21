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
	games      map[uuid.UUID]*GameServer
	Register   chan *GameServer
	Unregister chan uuid.UUID

	connection *rabbitmq.Conn
	logger     *logrus.Logger

	ActionChannel chan *GameActionRequest
}

func NewGameChannelServer(conn *rabbitmq.Conn, logger *logrus.Logger) *GameChannelServer {
	var server = &GameChannelServer{
		connection:    conn,
		logger:        logger,
		games:         make(map[uuid.UUID]*GameServer),
		Register:      make(chan *GameServer),
		Unregister:    make(chan uuid.UUID),
		ActionChannel: make(chan *GameActionRequest),
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
		case action := <-c.ActionChannel:
			c.action(action)
		}
	}
}

func (c *GameChannelServer) NewLiveGameChannel(instance *game.Instance) *GameServer {
	return NewGameServer(instance, c.connection, c.logger)
}

func (c *GameChannelServer) register(channel *GameServer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.games[channel.InstanceID] = channel
	channel.Start()
	c.logger.WithFields(logrus.Fields{"count": len(c.games)}).Info("game server registered")
}

func (c *GameChannelServer) unregister(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if channel, ok := c.games[id]; ok {
		close(channel.game.ActionChannel)
		close(channel.game.ChangeChannel)
	}
	delete(c.games, id)
	c.logger.WithFields(logrus.Fields{"count": len(c.games)}).Info("game server unregistered")
}

func (c *GameChannelServer) action(request *GameActionRequest) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if instance, exists := c.games[request.InstanceID]; exists {
		request.Action.Perform(instance.game)
	}
	c.logger.WithFields(logrus.Fields{"instance.id": request.InstanceID}).Info("game action received")
}
