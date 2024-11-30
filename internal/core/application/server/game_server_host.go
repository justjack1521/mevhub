package server

import (
	uuid "github.com/satori/go.uuid"
	"log/slog"
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

	logger *slog.Logger

	ActionChannel     chan *GameActionRequest
	GameServerFactory *GameServerFactory
}

func NewGameServerHost(logger *slog.Logger, factory *GameServerFactory) *GameServerHost {
	var server = &GameServerHost{
		logger:            logger,
		games:             make(map[uuid.UUID]*GameServer),
		Register:          make(chan *GameServer, 5),
		Unregister:        make(chan uuid.UUID, 5),
		ActionChannel:     make(chan *GameActionRequest, 5),
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

	if request.PartyID == uuid.Nil {
		return
	}

	instance, exists := h.games[request.PartyID]

	if exists == false {
		h.logger.With(
			slog.String("instance.id", request.PartyID.String()),
			slog.Group("action",
				slog.String("action.type", reflect.TypeOf(request.Action).String()),
			),
		).Info("game server action orphaned")
		return
	}

	instance.game.ActionChannel <- request.Action
	h.logger.With(
		slog.String("instance.id", request.PartyID.String()),
		slog.Group("action",
			slog.String("action.type", reflect.TypeOf(request.Action).String()),
		),
	).Info("game server action received")

}
