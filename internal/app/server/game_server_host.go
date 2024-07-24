package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/domain/game"
	"reflect"
	"sync"
)

type GameServerHost struct {
	mu         sync.Mutex
	games      map[uuid.UUID]*GameServer
	Register   chan *GameServer
	Unregister chan uuid.UUID

	connection *rabbitmq.Conn
	logger     *logrus.Logger

	ActionChannel chan *GameActionRequest
}

func NewGameServerHost(conn *rabbitmq.Conn, logger *logrus.Logger) *GameServerHost {
	var server = &GameServerHost{
		connection:    conn,
		logger:        logger,
		games:         make(map[uuid.UUID]*GameServer),
		Register:      make(chan *GameServer),
		Unregister:    make(chan uuid.UUID),
		ActionChannel: make(chan *GameActionRequest),
	}
	return server
}

func (c *GameServerHost) Run() {
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

func (c *GameServerHost) NewLiveGameChannel(instance *game.Instance) *GameServer {
	return NewGameServer(instance, c.connection, c.logger)
}

func (c *GameServerHost) register(channel *GameServer) {
	c.games[channel.InstanceID] = channel
	channel.Start()
	c.logger.WithFields(logrus.Fields{"count": len(c.games)}).Info("game server registered")
}

func (c *GameServerHost) unregister(id uuid.UUID) {
	if channel, ok := c.games[id]; ok {
		close(channel.game.ActionChannel)
		close(channel.game.ChangeChannel)
	}
	delete(c.games, id)
	c.logger.WithFields(logrus.Fields{"count": len(c.games)}).Info("game server unregistered")
}

func (c *GameServerHost) action(request *GameActionRequest) {
	if instance, exists := c.games[request.InstanceID]; exists {
		instance.game.ActionChannel <- request.Action
		c.logger.WithFields(logrus.Fields{
			"instance.id": request.InstanceID,
			"action.type": reflect.TypeOf(request.Action),
		}).Info("game action received")
	}
	c.logger.WithFields(logrus.Fields{
		"instance.id": request.InstanceID,
		"action.type": reflect.TypeOf(request.Action),
	}).Info("game action orphaned")
}
