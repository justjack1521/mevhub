package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"mevhub/internal/core/domain/game"
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
		case game.PlayerAddChange:
			s.HandlePlayerAddChange(actual)
		case game.PlayerRemoveChange:
			s.HandlePlayerRemoveChange(actual)
		case game.PlayerReadyChange:
			s.HandlePlayerReadyChange(actual)
		case game.PlayerEnqueueActionChange:
			s.HandlePlayerEnqueueActionChange(actual)
		case game.PlayerDequeueActionChange:
			s.HandlePlayerDequeueActionChange(actual)
		case game.PlayerLockActionChange:
			s.HandlePlayerLockActionChange(actual)
		case game.StateChange:
			s.HandleGameStateChange(actual)
		}
	}
}

func (s *GameServer) HandlePlayerRemoveChange(change game.PlayerRemoveChange) {
	delete(s.clients, change.PlayerID)
	var notification = &protomulti.GamePlayerRemoveNotification{
		GameId:      s.InstanceID.String(),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_PLAYER_REMOVE, notification)
}

func (s *GameServer) HandleGameStateChange(change game.StateChange) {

	s.logger.WithFields(logrus.Fields{
		"instance.id": s.InstanceID.String(),
		"state.name":  reflect.TypeOf(change.State),
	}).Info("game server state change")

	switch actual := change.State.(type) {
	case *game.EnemyTurnState:
		s.HandleEnemyTurnStateChange(actual)
	case *game.EndGameState:
		s.HandleEndGameStateChange(actual)
	}
}

func (s *GameServer) HandleEndGameStateChange(change *game.EndGameState) {

}

func (s *GameServer) HandleEnemyTurnStateChange(change *game.EnemyTurnState) {
	var queues = make([]*protomulti.ProtoGamePlayerActionQueue, len(change.QueuedActions))
	for i, queued := range change.QueuedActions {
		var queue = &protomulti.ProtoGamePlayerActionQueue{
			PlayerId: queued.PlayerID.String(),
			Actions:  make([]*protomulti.ProtoGameAction, len(queued.Actions)),
		}
		for j, action := range queued.Actions {
			queue.Actions[j] = &protomulti.ProtoGameAction{
				Action:    protomulti.GamePlayerActionType(action.ActionType),
				Target:    int32(action.Target),
				SlotIndex: int32(action.SlotIndex),
				ElementId: action.ElementID.String(),
			}
		}
		queues[i] = queue
	}
	var message = &protomulti.GameActionQueueConfirmNotification{
		PlayerActionQueue: queues,
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_QUEUE_CONFIRM, message)
}

func (s *GameServer) HandlePlayerEnqueueActionChange(change game.PlayerEnqueueActionChange) {
	var message = &protomulti.GameEnqueueActionNotification{
		GameId:      change.InstanceID.String(),
		PlayerIndex: int32(change.PartySlot),
		Action:      protomulti.GamePlayerActionType(change.ActionType),
		SlotIndex:   int32(change.SlotIndex),
		Target:      int32(change.Target),
		ElementId:   change.ElementID.String(),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION, message)
}

func (s *GameServer) HandlePlayerDequeueActionChange(change game.PlayerDequeueActionChange) {
	var message = &protomulti.GameDequeueActionNotification{
		GameId:      change.InstanceID.String(),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION, message)
}

func (s *GameServer) HandlePlayerLockActionChange(change game.PlayerLockActionChange) {
	var message = &protomulti.GameLockActionNotification{
		GameId:          change.InstanceID.String(),
		PlayerIndex:     int32(change.PartySlot),
		ActionLockIndex: int32(change.ActionLockIndex),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION, message)
}

func (s *GameServer) HandlePlayerAddChange(change game.PlayerAddChange) {
	s.clients[change.PlayerID] = &PlayerChannel{
		UserID:   change.UserID,
		PlayerID: change.PlayerID,
	}
}

func (s *GameServer) HandlePlayerReadyChange(change game.PlayerReadyChange) {
	var message = &protomulti.GamePlayerReadyNotification{
		GameId:      change.InstanceID.String(),
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
