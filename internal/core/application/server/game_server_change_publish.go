package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
)

// ErrUnhandledGameChange surfaces a change type with no dispatch case. Every
// change the domain can emit must be handled (or explicitly ignored) here; a
// silent fall-through previously masked the entire pipeline being dead.
var ErrUnhandledGameChange = func(change game.Change) error {
	return fmt.Errorf("unhandled game change %T (%s)", change, change.Identifier())
}

type changeMarshaller struct {
	playerRemove     translate.GamePlayerRemoveChangeMarshaller
	playerReady      translate.GamePlayerReadyChangeMarshaller
	playerLock       translate.GamePlayerLockActionChangeMarshaller
	playerEnqueue    translate.GamePlayerEnqueueActionChangeMarshaller
	playerDequeue    translate.GamePlayerDequeueActionChangeMarshaller
	playerDisconnect translate.GamePlayerDisconnectChangeMarshaller
	playerReconnect  translate.GamePlayerReconnectChangeMarshaller
	gameSync         translate.GameStateSyncChangeMarshaller
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
			playerRemove:     translate.NewGamePlayerRemoveChangeMarshaller(),
			playerReady:      translate.NewGamePlayerReadyChangeMarshaller(),
			playerLock:       translate.NewGamePlayerLockActionChangeMarshaller(),
			playerEnqueue:    translate.NewGamePlayerEnqueueActionChangeMarshaller(),
			playerDequeue:    translate.NewGamePlayerDequeueActionChangeMarshaller(),
			playerDisconnect: translate.NewGamePlayerDisconnectChangeMarshaller(),
			playerReconnect:  translate.NewGamePlayerReconnectChangeMarshaller(),
			gameSync:         translate.NewGameStateSyncChangeMarshaller(),
		},
	}
}

func (c *ChangeHandlerPublisher) Handle(svr *GameServer, change game.Change) error {

	if err := c.handler.Handle(svr, change); err != nil {
		return err
	}

	switch actual := change.(type) {
	case *action.PlayerAddChange:
		return c.HandlePlayerAddChange(svr, actual)
	case *action.PlayerRemoveChange:
		return c.HandlePlayerRemoveChange(svr, actual)
	case *action.PlayerReadyChange:
		return c.HandlePlayerReadyChange(svr, actual)
	case *action.PlayerEnqueueActionChange:
		return c.HandlePlayerEnqueueActionChange(svr, actual)
	case *action.PlayerDequeueActionChange:
		return c.HandlePlayerDequeueActionChange(svr, actual)
	case *action.PlayerLockActionChange:
		return c.HandlePlayerLockActionChange(svr, actual)
	case *action.StateChange:
		return c.HandleGameStateChange(svr, actual)
	case *action.HPConsensusChange:
		return c.HandleHPConsensusChange(svr, actual)
	case *action.GameStateSyncChange:
		return c.HandleGameStateSync(svr, actual)
	case *action.CatchUpChange:
		return c.HandleCatchUpChange(svr, actual)
	case *action.PlayerDisconnectChange:
		return c.HandlePlayerDisconnectChange(svr, actual)
	case *action.PlayerReconnectChange:
		return c.HandlePlayerReconnectChange(svr, actual)
	case *action.PartyAddChange:
		// Internal bookkeeping only; clients learn party composition from the
		// lobby, which is where it has always actually come from.
		return nil
	default:
		return ErrUnhandledGameChange(change)
	}
}

func (c *ChangeHandlerPublisher) HandleGameStateChange(svr *GameServer, change *action.StateChange) error {
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

	// QueuedActions is keyed by party index, which is not guaranteed dense
	// from zero, so the notification is built by append rather than indexing.
	var queues = make([]*protomulti.ProtoGamePartyActionQueue, 0, len(change.QueuedActions))

	for index, queued := range change.QueuedActions {
		var p = &protomulti.ProtoGamePartyActionQueue{
			PartyIndex:        int32(index),
			PlayerActionQueue: make([]*protomulti.ProtoGamePlayerActionQueue, 0, len(queued)),
		}
		for _, q := range queued {
			if q == nil {
				continue
			}
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
			p.PlayerActionQueue = append(p.PlayerActionQueue, player)
		}
		queues = append(queues, p)
	}

	var message = &protomulti.GameActionQueueConfirmNotification{
		PartyActionQueues: queues,
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_QUEUE_CONFIRM, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerLockActionChange(svr *GameServer, change *action.PlayerLockActionChange) error {
	message, err := c.marshaller.playerLock.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerDequeueActionChange(svr *GameServer, change *action.PlayerDequeueActionChange) error {
	message, err := c.marshaller.playerDequeue.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerEnqueueActionChange(svr *GameServer, change *action.PlayerEnqueueActionChange) error {
	message, err := c.marshaller.playerEnqueue.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerAddChange(svr *GameServer, change *action.PlayerAddChange) error {
	return nil
}

// HandlePlayerDisconnectChange broadcasts a disconnect notification.
func (c *ChangeHandlerPublisher) HandlePlayerDisconnectChange(svr *GameServer, change *action.PlayerDisconnectChange) error {
	message, err := c.marshaller.playerDisconnect.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_DISCONNECT, message)
}

// HandlePlayerReconnectChange broadcasts a reconnect notification.
func (c *ChangeHandlerPublisher) HandlePlayerReconnectChange(svr *GameServer, change *action.PlayerReconnectChange) error {
	message, err := c.marshaller.playerReconnect.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_RECONNECT, message)
}

func (c *ChangeHandlerPublisher) HandlePlayerRemoveChange(svr *GameServer, change *action.PlayerRemoveChange) error {
	message, err := c.marshaller.playerRemove.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_REMOVE, message)
}

// HandleGameStateSync delivers a full live-state snapshot to a single
// reconnecting player rather than broadcasting it. Nothing emits this change
// any more — notification replay took its place — but the path is kept intact
// and reversible.
func (c *ChangeHandlerPublisher) HandleGameStateSync(svr *GameServer, change *action.GameStateSyncChange) error {
	message, err := c.marshaller.gameSync.Marshall(change)
	if err != nil {
		return err
	}
	svr.mu.RLock()
	client, ok := svr.clients[change.TargetPlayerID]
	svr.mu.RUnlock()
	if !ok {
		return nil
	}
	return c.publishTo(client, protomulti.MultiGameNotificationType_GAME_NOTIFY_GAME_SYNC, message)
}

// HandleCatchUpChange replays the log to the one player who asked for it. The
// replayed notifications keep their original kind, payload and ordinal, so the
// client cannot tell them from live traffic and can dedupe against what it has
// already applied.
func (c *ChangeHandlerPublisher) HandleCatchUpChange(svr *GameServer, change *action.CatchUpChange) error {

	// Sampled before the send so the caller learns the range even if it has no
	// delivery path left; the log holds everything published up to this change,
	// because appends happen on this goroutine.
	entries, from, to := svr.notifications.Replay(change.FromSequence)
	change.Report(action.CatchUpResult{
		FromSequence:    from,
		ToSequence:      to,
		TurnRemainingMs: change.TurnRemainingMs,
	})

	svr.mu.RLock()
	client, ok := svr.clients[change.TargetPlayerID]
	svr.mu.RUnlock()
	if !ok {
		return nil
	}

	// One client's failure must not abandon the rest of its own backlog: a gap
	// mid-replay is exactly what the client would have to ask about again.
	var errs []error
	for _, entry := range entries {
		var notification = &protocommon.Notification{
			Service:  protocommon.ServiceKey_MULTI,
			Type:     int32(entry.kind),
			Data:     entry.payload,
			Sequence: entry.sequence,
		}
		if err := c.publisher.Publish(context.Background(), client, notification); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)

}

// publishTo delivers to a single client and is never logged for replay: its
// traffic is targeted at one player and must not become replayable by another.
// It therefore goes out unsequenced — taking an ordinal from the broadcast
// stream would put a hole in that stream for everybody else.
func (c *ChangeHandlerPublisher) publishTo(client *PlayerChannel, operation protomulti.MultiGameNotificationType, message Notification) error {
	bytes, err := message.MarshallBinary()
	if err != nil {
		return err
	}
	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(operation),
		Data:    bytes,
	}
	return c.publisher.Publish(context.Background(), client, notification)
}

func (c *ChangeHandlerPublisher) HandlePlayerReadyChange(svr *GameServer, change *action.PlayerReadyChange) error {
	message, err := c.marshaller.playerReady.Marshall(change)
	if err != nil {
		return err
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_READY, message)
}

func (c *ChangeHandlerPublisher) HandleHPConsensusChange(svr *GameServer, change *action.HPConsensusChange) error {
	enemies := make([]*protomulti.ProtoGameEnemyHP, len(change.Enemies))
	for i, e := range change.Enemies {
		enemies[i] = &protomulti.ProtoGameEnemyHP{
			EnemyIndex: int32(e.EnemyIndex),
			Hp:         int32(e.HP),
		}
	}
	message := &protomulti.GameHPSyncNotification{
		GameId:  svr.InstanceID.String(),
		Enemies: enemies,
	}
	return c.publish(svr, protomulti.MultiGameNotificationType_GAME_NOTIFY_HP_CONSENSUS, message)
}

func (c *ChangeHandlerPublisher) publish(svr *GameServer, operation protomulti.MultiGameNotificationType, message Notification) error {

	bytes, err := message.MarshallBinary()
	if err != nil {
		return err
	}

	// Log before the fan-out, not after: a publish that errors below is then
	// still replayable, which is the whole point of keeping the bytes. The log
	// also assigns the ordinal, so what a client sees is numbered gaplessly.
	var sequence = svr.notifications.Append(operation, bytes)

	var notification = &protocommon.Notification{
		Service:  protocommon.ServiceKey_MULTI,
		Type:     int32(operation),
		Data:     bytes,
		Sequence: sequence,
	}

	// Snapshot under the lock, publish outside it: the client map is mutated
	// from other goroutines and network I/O must not hold the lock. One
	// client's failure must not starve the rest of the broadcast.
	svr.mu.RLock()
	var clients = make([]*PlayerChannel, 0, len(svr.clients))
	for _, client := range svr.clients {
		clients = append(clients, client)
	}
	svr.mu.RUnlock()

	var errs []error
	for _, client := range clients {
		if err := c.publisher.Publish(context.Background(), client, notification); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
