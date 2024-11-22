package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log/slog"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
)

type GameServerFactory struct {
	actions []GameServerFactoryBuildAction
}

func NewGameServerFactory(actions []GameServerFactoryBuildAction) *GameServerFactory {
	return &GameServerFactory{actions: actions}
}

func (f *GameServerFactory) Create(instance *game.Instance, notifier NotificationPublisher) *GameServer {
	var svr = &GameServer{
		InstanceID:   instance.SysID,
		game:         game.NewLiveGameInstance(instance),
		clients:      make(map[uuid.UUID]*PlayerChannel),
		publisher:    notifier,
		logger:       logrus.New(),
		ErrorHandler: ErrorHandlerDefault{},
	}
	for _, a := range f.actions {
		a(svr)
	}
	svr.game.State = action.NewPendingState(svr.game)
	return svr
}

type GameServerFactoryBuildAction func(svr *GameServer)

func GameServerFactoryLoggingBuildAction(logger *slog.Logger) GameServerFactoryBuildAction {
	return func(svr *GameServer) {
		svr.ErrorHandler = NewErrorLoggingDecorator(logger, svr.ErrorHandler)
	}
}
