package server

import (
	uuid "github.com/satori/go.uuid"
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

func (f *GameServerFactory) Create(instance *game.Instance) *GameServer {
	var svr = &GameServer{
		InstanceID:    instance.SysID,
		game:          game.NewLiveGameInstance(instance),
		clients:       make(map[uuid.UUID]*PlayerChannel),
		ChangeHandler: NewChangeHandlerDefault(),
		ErrorHandler:  NewErrorHandlerDefault(),
	}
	for _, a := range f.actions {
		a(svr)
	}
	svr.game.State = action.NewPendingState(svr.game)
	return svr
}

type GameServerFactoryBuildAction func(svr *GameServer)

func GameServerFactoryPublisherBuildAction(publisher NotificationPublisher) GameServerFactoryBuildAction {
	return func(svr *GameServer) {
		svr.ChangeHandler = NewChangeHandlerPublisher(publisher, svr.ChangeHandler)
	}
}

func GameServerFactoryLoggingBuildAction(logger *slog.Logger) GameServerFactoryBuildAction {
	return func(svr *GameServer) {
		svr.ErrorHandler = NewErrorLoggingDecorator(logger, svr.ErrorHandler)
	}
}
