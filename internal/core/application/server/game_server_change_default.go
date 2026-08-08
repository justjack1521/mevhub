package server

import (
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
	"time"
)

type ChangeHandlerDefault struct {
}

func NewChangeHandlerDefault() *ChangeHandlerDefault {
	return &ChangeHandlerDefault{}
}

func (c *ChangeHandlerDefault) Handle(svr *GameServer, change game.Change) error {
	// This handler maintains the client registry only; change types it does
	// not care about intentionally fall through. Completeness is enforced by
	// ChangeHandlerPublisher's default arm.
	switch actual := change.(type) {
	case *action.PlayerAddChange:
		return c.HandlePlayerAddChance(svr, actual)
	case *action.PlayerRemoveChange:
		return c.HandlePlayerRemoveChance(svr, actual)
	case *action.PlayerDisconnectChange:
		return c.HandlePlayerDisconnectChange(svr, actual)
	case *action.PlayerReconnectChange:
		return c.HandlePlayerReconnectChange(svr, actual)
	}
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerAddChance(svr *GameServer, change *action.PlayerAddChange) error {
	svr.mu.Lock()
	svr.clients[change.PlayerID] = &PlayerChannel{
		UserID:   change.UserID,
		PlayerID: change.PlayerID,
	}
	svr.mu.Unlock()
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerRemoveChance(svr *GameServer, change *action.PlayerRemoveChange) error {
	svr.mu.Lock()
	delete(svr.clients, change.PlayerID)
	svr.mu.Unlock()
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerDisconnectChange(svr *GameServer, change *action.PlayerDisconnectChange) error {
	svr.mu.Lock()
	if ch, ok := svr.clients[change.PlayerID]; ok {
		now := time.Now().UTC()
		ch.DisconnectedAt = &now
		ch.timedOut = false
	}
	svr.mu.Unlock()
	return nil
}

func (c *ChangeHandlerDefault) HandlePlayerReconnectChange(svr *GameServer, change *action.PlayerReconnectChange) error {
	svr.mu.Lock()
	if ch, ok := svr.clients[change.PlayerID]; ok {
		ch.DisconnectedAt = nil
		ch.timedOut = false
	}
	svr.mu.Unlock()
	return nil
}
