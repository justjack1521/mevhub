package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
)

type changeMarshaller struct {
	playerRemove  translate.GamePlayerRemoveChangeMarshaller
	playerReady   translate.GamePlayerReadyChangeMarshaller
	playerLock    translate.GamePlayerLockActionChangeMarshaller
	playerEnqueue translate.GamePlayerEnqueueActionChangeMarshaller
	playerDequeue translate.GamePlayerDequeueActionChangeMarshaller
}

type ChangeHandlerPublisher struct {
	handler        ChangeHandler
	publisher      NotificationPublisher
	eventPublisher *mevent.Publisher
	marshaller     changeMarshaller
}

func NewChangeHandlerPublisher(publisher NotificationPublisher, eventPublisher *mevent.Publisher, handler ChangeHandler) *ChangeHandlerPublisher {
	return &ChangeHandlerPublisher{
		publisher:      publisher,
		eventPublisher: eventPublisher,
		handler:        handler,
		marshaller: changeMarshaller{
			playerRemove:  translate.NewGamePlayerRemoveChangeMarshaller(),
			playerReady:   translate.NewGamePlayerReadyChangeMarshaller(),
			playerLock:    translate.NewGamePlayerLockActionChangeMarshaller(),
			playerEnqueue: translate.NewGamePlayerEnqueueActionChangeMarshaller(),
			playerDequeue: translate.NewGamePlayerDequeueActionChangeMarshaller(),
		},
	}
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
	case *action.PlayerTurnState:
		return c.HandlePlayerTurnStateChange(svr, actual)
	case *action.EnemyTurnState:
		return c.HandleEnemyTurnStateChange(svr, actual)
	case *action.EndGameState:
		return c.HandleEndGameStateChange(svr, actual)
	}
	return nil
}

func (c *ChangeHandlerPublisher) HandleEndGameStateChange(svr *GameServer, _ *action.EndGameState) error {
	var message = &protomulti.GameEndNotification{
		GameId: svr.InstanceID.String(),
	}
	if err := c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_END, message); err != nil {
		return err
	}
	c.eventPublisher.Notify(game.NewInstanceDeletedEvent(context.Background(), svr.InstanceID))
	return nil
}

func (c *ChangeHandlerPublisher) HandlePlayerTurnStateChange(svr *GameServer, _ *action.PlayerTurnState) error {
	var message = &protomulti.GameReadyNotification{
		GameId: svr.InstanceID.String(),
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_READY, message)
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
				var act = &protomulti.ProtoGameAction{
					Action:    protomulti.GamePlayerActionType(a.ActionType),
					Target:    int32(a.Target),
					SlotIndex: int32(a.SlotIndex),
					ElementId: a.ElementID.String(),
				}
				player.Actions[k] = act
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
	message, err := c.marshaller.playerLock.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerDequeueActionChange(svr *GameServer, change action.PlayerDequeueActionChange) error {
	message, err := c.marshaller.playerDequeue.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerEnqueueActionChange(svr *GameServer, change action.PlayerEnqueueActionChange) error {
	message, err := c.marshaller.playerEnqueue.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerAddChange(svr *GameServer, change action.PlayerAddChange) error {
	return nil
}

func (c *ChangeHandlerPublisher) HandlePlayerRemoveChange(svr *GameServer, change action.PlayerRemoveChange) error {
	message, err := c.marshaller.playerRemove.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_REMOVE, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerReadyChange(svr *GameServer, change action.PlayerReadyChange) error {
	message, err := c.marshaller.playerReady.Marshall(change)
	if err != nil {
		return err
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
