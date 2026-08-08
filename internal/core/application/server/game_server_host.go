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

const gameServerHostReapCheckPeriod = time.Second * 30

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

	// Unregister directly rather than sending to h.Unregister: this runs on
	// the goroutine that drains that channel, so a send here deadlocks the
	// whole host once the buffer fills.
	var ended []uuid.UUID
	var timedOut []game.PlayerTimedOutEvent

	for id, instance := range h.games {
		if instance.game.HasEnded() {
			ended = append(ended, id)
			continue
		}
		for _, ch := range instance.ClaimExpiredClients(ClientTimeoutPeriod) {
			timedOut = append(timedOut, game.NewPlayerTimedOutEvent(
				context.Background(), instance.InstanceID, ch.UserID, ch.PlayerID,
			))
		}
	}

	for _, id := range ended {
		h.unregister(id)
	}

	// One goroutine per reap pass, not per client: subscribers do Redis I/O
	// that must not block the host, but goroutine growth stays bounded.
	if len(timedOut) > 0 {
		go func(events []game.PlayerTimedOutEvent) {
			for _, event := range events {
				h.eventPublisher.Notify(event)
			}
		}(timedOut)
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
		// Signal-based shutdown: closing the data channels would panic any
		// goroutine still sending on them.
		channel.Stop()
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

	select {
	case instance.game.ActionChannel <- request.Action:
	case <-instance.game.Done():
		// The game shut down while routing; nothing left to receive it.
		return
	}
	h.logger.With(
		slog.String("instance.id", request.GameID.String()),
		slog.Group("action",
			slog.String("action.type", reflect.TypeOf(request.Action).String()),
		),
	).Info("game server action received")

}
