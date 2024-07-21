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
		case game.PlayerReadyChange:
			s.HandlePlayerReadyChange(actual)
		case game.PlayerAddChange:
			s.HandlePlayerAddChange(actual)
		}
	}
}

func (s *GameServer) HandlePlayerAddChange(request game.PlayerAddChange) {
	s.mu.Lock()
	s.clients[request.PlayerID] = &PlayerChannel{
		UserID:   request.UserID,
		PlayerID: request.PlayerID,
	}
	s.mu.Unlock()
}

func (s *GameServer) HandlePlayerReadyChange(change game.PlayerReadyChange) {
	var message = &protomulti.GamePlayerReadyNotification{
		GameId:      change.InstanceID.String(),
		PlayerIndex: int32(change.PartySlot),
	}
	s.Publish(protomulti.MultiGameNotificationType_GAME_PLAYER_READY, message)
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

	s.mu.RLock()
	for _, client := range s.clients {
		if err := s.publisher.Publish(context.Background(), msg, client.UserID, client.PlayerID, rabbitmv.ClientNotification); err != nil {
			return
		}
	}
	s.mu.RUnlock()

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
