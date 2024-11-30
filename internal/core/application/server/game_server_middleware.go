package server

import (
	"mevhub/internal/core/domain/game"
)

type ChangeHandler interface {
	Handle(svr *GameServer, change game.Change) error
}

type ErrorHandler interface {
	Handle(svr *GameServer, err error)
}
