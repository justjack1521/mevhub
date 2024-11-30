package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
)

type ChangeHandlerPublisher struct {
	handler   ChangeHandler
	publisher NotificationPublisher
}

func NewChangeHandlerPublisher(publisher NotificationPublisher, handler ChangeHandler) *ChangeHandlerPublisher {
	return &ChangeHandlerPublisher{publisher: publisher, handler: handler}
}

func (c *ChangeHandlerPublisher) Handle(svr *GameServer, change game.Change) error {

	if err := c.handler.Handle(svr, change); err != nil {
		return err
	}

	switch actual := change.(type) {
	case action.PlayerAddChange:
		return c.HandlePlayerAddChange(svr, actual)
	case action.PlayerRemoveChange:
		return c.HandlePlayerRemoveChange(svr, actual)
	case action.PlayerReadyChange:
		return c.HandlePlayerReadyChange(svr, actual)
	case action.PlayerEnqueueActionChange:
		return c.HandlePlayerEnqueueActionChange(svr, actual)
	case action.PlayerDequeueActionChange:
		return c.HandlePlayerDequeueActionChange(svr, actual)
	case action.PlayerLockActionChange:
		return c.HandlePlayerLockActionChange(svr, actual)
	case action.StateChange:
		return c.HandleGameStateChange(svr, actual)
	}
	return nil
}

func (c *ChangeHandlerPublisher) HandleGameStateChange(svr *GameServer, change action.StateChange) error {
	switch actual := change.State.(type) {
	case *action.EnemyTurnState:
		return c.HandleEnemyTurnStateChange(svr, actual)
	}
	return nil
}

func (c *ChangeHandlerPublisher) HandleEnemyTurnStateChange(svr *GameServer, change *action.EnemyTurnState) error {

	var queues = make([]*protomulti.ProtoGamePartyActionQueue, len(change.QueuedActions))

	for index, queued := range change.QueuedActions {
		var p = &protomulti.ProtoGamePartyActionQueue{
			PartyIndex:        int32(index),
			PlayerActionQueue: make([]*protomulti.ProtoGamePlayerActionQueue, len(queued)),
		}
		for i, q := range queued {
			var player = &protomulti.ProtoGamePlayerActionQueue{
				PlayerId: q.PlayerID.String(),
				Actions:  make([]*protomulti.ProtoGameAction, len(q.Actions)),
			}
			for k, a := range q.Actions {
				var action = &protomulti.ProtoGameAction{
					Action:    protomulti.GamePlayerActionType(a.ActionType),
					Target:    int32(a.Target),
					SlotIndex: int32(a.SlotIndex),
					ElementId: a.ElementID.String(),
				}
				player.Actions[k] = action
			}
			p.PlayerActionQueue[i] = player
		}
	}

	var message = &protomulti.GameActionQueueConfirmNotification{
		PartyActionQueues: queues,
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_QUEUE_CONFIRM, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerLockActionChange(svr *GameServer, change action.PlayerLockActionChange) error {
	var message = &protomulti.GameLockActionNotification{
		GameId:          change.InstanceID.String(),
		PartyIndex:      int32(change.PartyIndex),
		PlayerIndex:     int32(change.PartySlot),
		ActionLockIndex: int32(change.ActionLockIndex),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerDequeueActionChange(svr *GameServer, change action.PlayerDequeueActionChange) error {
	var message = &protomulti.GameDequeueActionNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerEnqueueActionChange(svr *GameServer, change action.PlayerEnqueueActionChange) error {
	var message = &protomulti.GameEnqueueActionNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
		Action:      protomulti.GamePlayerActionType(change.ActionType),
		SlotIndex:   int32(change.SlotIndex),
		Target:      int32(change.Target),
		ElementId:   change.ElementID.String(),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerAddChange(svr *GameServer, change action.PlayerAddChange) error {
	return nil
}

func (c *ChangeHandlerPublisher) HandlePlayerRemoveChange(svr *GameServer, change action.PlayerRemoveChange) error {
	var notification = &protomulti.GamePlayerRemoveNotification{
		GameId:      svr.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_REMOVE, notification)
}

func (c *ChangeHandlerPublisher) HandlePlayerReadyChange(svr *GameServer, change action.PlayerReadyChange) error {
	var message = &protomulti.GamePlayerReadyNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_READY, message)
}

func (c *ChangeHandlerPublisher) publish(svr *GameServer, operation protomulti.MultiGameNotificationType, message Notification) error {

	bytes, err := message.MarshallBinary()
	if err != nil {
		return err
	}

	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(operation),
		Data:    bytes,
	}

	for _, client := range svr.clients {
		if err := c.publisher.Publish(context.Background(), client, notification); err != nil {
			return err
		}
	}

	return nil
}
