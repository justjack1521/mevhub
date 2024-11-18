package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/core/domain/game"
	"reflect"
	"sync"
	"time"
)

const gameServerHostReapCheckPeriod = time.Minute * 3

type GameServerHost struct {
	mu         sync.Mutex
	games      map[uuid.UUID]*GameServer
	Register   chan *GameServer
	Unregister chan uuid.UUID

	connection *rabbitmq.Conn
	logger     *logrus.Logger

	ActionChannel     chan *GameActionRequest
	GameServerFactory *GameServerFactory
}

func NewGameServerHost(conn *rabbitmq.Conn, logger *logrus.Logger, factory *GameServerFactory) *GameServerHost {
	var server = &GameServerHost{
		connection:        conn,
		logger:            logger,
		games:             make(map[uuid.UUID]*GameServer),
		Register:          make(chan *GameServer),
		Unregister:        make(chan uuid.UUID),
		ActionChannel:     make(chan *GameActionRequest),
		GameServerFactory: factory,
	}
	return server
}

func (h *GameServerHost) Run() {

	var ticker = time.NewTicker(gameServerHostReapCheckPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case c := <-ticker.C:
			h.tick(c)
		case instance := <-h.Register:
			h.register(instance)
		case id := <-h.Unregister:
			h.unregister(id)
		case action := <-h.ActionChannel:
			h.action(action)
		}
	}
}

func (h *GameServerHost) NewLiveGameChannel(instance *game.Instance) *GameServer {
	return h.GameServerFactory.Create(instance, NewGameServerRabbitMQNotifier(h.connection))
}

func (h *GameServerHost) tick(t time.Time) {
	for id, instance := range h.games {
		if instance.game.Ended {
			h.Unregister <- id
		}
	}
}

func (h *GameServerHost) register(channel *GameServer) {
	h.games[channel.InstanceID] = channel
	channel.Start()
	h.logger.WithFields(logrus.Fields{"count": len(h.games)}).Info("game server registered")
}

func (h *GameServerHost) unregister(id uuid.UUID) {
	if channel, ok := h.games[id]; ok {
		close(channel.game.ActionChannel)
		close(channel.game.ChangeChannel)
		close(channel.game.ErrorChannel)
	}
	delete(h.games, id)
	h.logger.WithFields(logrus.Fields{"count": len(h.games)}).Info("game server unregistered")
}

func (h *GameServerHost) action(request *GameActionRequest) {

	if request.InstanceID == uuid.Nil {
		return
	}

	instance, exists := h.games[request.InstanceID]

	if exists == false {
		h.logger.WithFields(logrus.Fields{
			"instance.id": request.InstanceID,
			"action.type": reflect.TypeOf(request.Action),
		}).Info("game server action orphaned")
		return
	}

	instance.game.ActionChannel <- request.Action
	h.logger.WithFields(logrus.Fields{
		"instance.id": request.InstanceID,
		"action.type": reflect.TypeOf(request.Action),
	}).Info("game server action received")

}
