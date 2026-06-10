package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"log/slog"
	"mevhub/internal/core/domain/game"
	"reflect"
	"time"
)

const gameServerHostReapCheckPeriod = time.Minute * 3

type GameServerHost struct {
	games      map[uuid.UUID]*GameServer
	Register   chan *GameServer
	Unregister chan uuid.UUID

	logger         *slog.Logger
	eventPublisher *mevent.Publisher

	ActionChannel     chan *GameActionRequest
	GameServerFactory *GameServerFactory
}

func NewGameServerHost(logger *slog.Logger, factory *GameServerFactory, publisher *mevent.Publisher) *GameServerHost {
	var server = &GameServerHost{
		logger:            logger,
		eventPublisher:    publisher,
		games:             make(map[uuid.UUID]*GameServer),
		Register:          make(chan *GameServer, 5),
		Unregister:        make(chan uuid.UUID, 5),
		ActionChannel:     make(chan *GameActionRequest, 64),
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
	return h.GameServerFactory.Create(instance)
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
	h.logger.With(slog.Int("count", len(h.games))).Info("game server registered")
	go h.eventPublisher.Notify(game.NewInstanceRegisteredEvent(context.Background(), channel.InstanceID))
}

func (h *GameServerHost) unregister(id uuid.UUID) {
	if channel, ok := h.games[id]; ok {
		close(channel.game.ActionChannel)
		close(channel.game.ChangeChannel)
		close(channel.game.ErrorChannel)
	}
	delete(h.games, id)
	h.logger.With(slog.Int("count", len(h.games))).Info("game server unregistered")
}

func (h *GameServerHost) action(request *GameActionRequest) {

	if request.GameID == uuid.Nil {
		return
	}

	instance, exists := h.games[request.GameID]

	if exists == false {
		h.logger.With(
			slog.String("instance.id", request.GameID.String()),
			slog.Group("action",
				slog.String("action.type", reflect.TypeOf(request.Action).String()),
			),
		).Info("game server action orphaned")
		return
	}

	instance.game.ActionChannel <- request.Action
	h.logger.With(
		slog.String("instance.id", request.GameID.String()),
		slog.Group("action",
			slog.String("action.type", reflect.TypeOf(request.Action).String()),
		),
	).Info("game server action received")

}
