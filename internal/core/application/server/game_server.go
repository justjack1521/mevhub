package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
	"reflect"
	"sync"
	"time"
)

const (
	ClientTimeoutPeriod = time.Minute * 3
)

type NotificationPublisher interface {
	Publish(ctx context.Context, player *PlayerChannel, notification Notification) error
}

type Notification interface {
	MarshallBinary() ([]byte, error)
}

type GameServer struct {
	InstanceID   uuid.UUID
	game         *game.LiveGameInstance
	mu           sync.RWMutex
	clients      map[uuid.UUID]*PlayerChannel
	publisher    NotificationPublisher
	logger       *logrus.Logger
	ErrorHandler ErrorHandler
	errorCount   int
}

func (s *GameServer) Start() {
	go s.WatchChanges()
	go s.WatchErrors()
	go s.game.WatchActions()
	go s.game.Tick()
}

func (s *GameServer) WatchErrors() {
	for {
		err := <-s.game.ErrorChannel
		s.ErrorHandler.Handle(s, err)
	}
}

func (s *GameServer) WatchChanges() {
	for {

		change := <-s.game.ChangeChannel

		s.logger.WithFields(logrus.Fields{
			"instance.id":       s.InstanceID.String(),
			"change.identifier": change.Identifier(),
		}).Info("game server change received")

		switch actual := change.(type) {
		case action.PlayerAddChange:
			s.HandlePlayerAddChange(actual)
		case action.PlayerRemoveChange:
			s.HandlePlayerRemoveChange(actual)
		case action.PlayerReadyChange:
			s.HandlePlayerReadyChange(actual)
		case action.PlayerEnqueueActionChange:
			s.HandlePlayerEnqueueActionChange(actual)
		case action.PlayerDequeueActionChange:
			s.HandlePlayerDequeueActionChange(actual)
		case action.PlayerLockActionChange:
			s.HandlePlayerLockActionChange(actual)
		case action.StateChange:
			s.HandleGameStateChange(actual)
		}
	}
}

func (s *GameServer) HandlePlayerRemoveChange(change action.PlayerRemoveChange) {
	delete(s.clients, change.PlayerID)
	var notification = &protomulti.GamePlayerRemoveNotification{
		GameId:      s.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_REMOVE, notification)
}

func (s *GameServer) HandleGameStateChange(change action.StateChange) {

	s.logger.WithFields(logrus.Fields{
		"instance.id": s.InstanceID.String(),
		"state.name":  reflect.TypeOf(change.State),
	}).Info("game server state change")

	switch actual := change.State.(type) {
	case *action.EnemyTurnState:
		s.HandleEnemyTurnStateChange(actual)
	case *action.EndGameState:
		s.HandleEndGameStateChange(actual)
	}
}

func (s *GameServer) HandleEndGameStateChange(change *action.EndGameState) {

}

func (s *GameServer) HandleEnemyTurnStateChange(change *action.EnemyTurnState) {

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
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_QUEUE_CONFIRM, message)
}

func (s *GameServer) HandlePlayerEnqueueActionChange(change action.PlayerEnqueueActionChange) {
	var message = &protomulti.GameEnqueueActionNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
		Action:      protomulti.GamePlayerActionType(change.ActionType),
		SlotIndex:   int32(change.SlotIndex),
		Target:      int32(change.Target),
		ElementId:   change.ElementID.String(),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION, message)
}

func (s *GameServer) HandlePlayerDequeueActionChange(change action.PlayerDequeueActionChange) {
	var message = &protomulti.GameDequeueActionNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION, message)
}

func (s *GameServer) HandlePlayerLockActionChange(change action.PlayerLockActionChange) {
	var message = &protomulti.GameLockActionNotification{
		GameId:          change.InstanceID.String(),
		PartyIndex:      int32(change.PartyIndex),
		PlayerIndex:     int32(change.PartySlot),
		ActionLockIndex: int32(change.ActionLockIndex),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION, message)
}

func (s *GameServer) HandlePlayerAddChange(change action.PlayerAddChange) {
	s.clients[change.PlayerID] = &PlayerChannel{
		UserID:   change.UserID,
		PlayerID: change.PlayerID,
	}
}

func (s *GameServer) HandlePlayerReadyChange(change action.PlayerReadyChange) {
	var message = &protomulti.GamePlayerReadyNotification{
		GameId:      change.InstanceID.String(),
		PartyIndex:  int32(change.PartyIndex),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_READY, message)
}

func (s *GameServer) Publish(operation protomulti.MultiGameNotificationType, message Notification) {

	bytes, err := message.MarshallBinary()
	if err != nil {
		return
	}

	s.logger.WithFields(logrus.Fields{
		"operation":    operation,
		"length":       len(bytes),
		"instance.id":  s.InstanceID.String(),
		"player.count": len(s.clients),
	}).Info("game server dispatching notification")

	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(operation),
		Data:    bytes,
	}

	for _, client := range s.clients {
		if err := s.publisher.Publish(context.Background(), client, notification); err != nil {
			return
		}
	}

}

type GameActionRequest struct {
	GameID  uuid.UUID
	PartyID uuid.UUID
	Action  game.Action
}

type PlayerAddRequest struct {
	PartyID   uuid.UUID
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	PartySlot int
}

type PlayerRemoveRequest struct {
	PlayerID uuid.UUID
}
