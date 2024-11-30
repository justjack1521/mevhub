package server

import (
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
)

type ChangeHandlerDefault struct {
}

func NewChangeHandlerDefault() *ChangeHandlerDefault {
	return &ChangeHandlerDefault{}
}

func (c *ChangeHandlerDefault) Handle(svr *GameServer, change game.Change) error {
	switch actual := change.(type) {
	case action.PlayerAddChange:
		svr.clients[actual.PlayerID] = &PlayerChannel{
			UserID:   actual.UserID,
			PlayerID: actual.PlayerID,
		}
	case action.PlayerRemoveChange:
		delete(svr.clients, actual.PlayerID)
	}
	return nil
}
