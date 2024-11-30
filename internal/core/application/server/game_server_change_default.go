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
		return c.HandlePlayerAddChance(svr, actual)
	case action.PlayerRemoveChange:
		return c.HandlePlayerRemoveChance(svr, actual)
	}
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerAddChance(svr *GameServer, change action.PlayerAddChange) error {
	svr.clients[change.PlayerID] = &PlayerChannel{
		UserID:   change.UserID,
		PlayerID: change.PlayerID,
	}
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerRemoveChance(svr *GameServer, change action.PlayerRemoveChange) error {
	delete(svr.clients, change.PlayerID)
	return nil
}
