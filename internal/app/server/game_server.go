package server

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/rabbitmv"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/domain/game"
	"sync"
)

type Notification interface {
	MarshallBinary() ([]byte, error)
}

type GameServer struct {
	InstanceID uuid.UUID
	game       *game.LiveGameInstance
	mu         sync.RWMutex
	clients    map[uuid.UUID]*PlayerChannel
	publisher  *rabbitmv.StandardPublisher
}

func (s *GameServer) Start() {
	go s.WatchChanges()
	go s.game.WatchActions()
	go s.game.Tick()
}

func NewGameServer(instance *game.Instance, conn *rabbitmq.Conn, logger *logrus.Logger) *GameServer {
	return &GameServer{
		InstanceID: instance.SysID,
		game:       game.NewLiveGameInstance(),
		clients:    make(map[uuid.UUID]*PlayerChannel),
		publisher:  rabbitmv.NewClientPublisher(conn, rabbitmq.WithPublisherOptionsLogger(logger)),
	}
}

func (s *GameServer) WatchChanges() {
	for {
		change := <-s.game.ChangeChannel
		switch actual := change.(type) {
		case game.PlayerAddChange:
			s.HandlePlayerAddChange(actual)
		case game.PlayerReadyChange:
			s.HandlePlayerReadyChange(actual)
		case game.PlayerEnqueueActionChange:
			s.HandlePlayerEnqueueActionChange(actual)
		case game.PlayerDequeueActionChange:
			s.HandlePlayerDequeueActionChange(actual)
		case game.PlayerLockActionChange:
			s.HandlePlayerLockActionChange(actual)
		}
	}
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

	var notification = &protocommon.Notification{
		Service: protocommon.ServiceKey_MULTI,
		Type:    int32(operation),
		Data:    bytes,
	}

	msg, err := notification.MarshallBinary()
	if err != nil {
		return
	}

	for _, client := range s.clients {
		if err := s.publisher.Publish(context.Background(), msg, client.UserID, client.PlayerID, rabbitmv.ClientNotification); err != nil {
			return
		}
	}

}

type GameActionRequest struct {
	InstanceID uuid.UUID
	Action     game.Action
}

type PlayerAddRequest struct {
	InstanceID uuid.UUID
	UserID     uuid.UUID
	PlayerID   uuid.UUID
	PartySlot  int
}

type PlayerRemoveRequest struct {
	PlayerID uuid.UUID
}
