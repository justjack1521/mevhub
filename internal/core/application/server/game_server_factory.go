package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"mevhub/internal/core/domain/game"
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
		ErrorHandler: ErrorHandlerDefault{},
	}
	for _, action := range f.actions {
		action(svr)
	}
	return svr
}

type GameServerFactoryBuildAction func(svr *GameServer)

func GameServerFactoryLoggingBuildAction(logger *logrus.Logger) GameServerFactoryBuildAction {
	return func(svr *GameServer) {
		svr.ErrorHandler = NewErrorLoggingDecorator(logger, svr.ErrorHandler)
	}
}
