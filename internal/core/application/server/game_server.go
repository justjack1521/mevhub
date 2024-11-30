package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
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
	InstanceID    uuid.UUID
	game          *game.LiveGameInstance
	mu            sync.RWMutex
	clients       map[uuid.UUID]*PlayerChannel
	ChangeHandler ChangeHandler
	ErrorHandler  ErrorHandler
	errorCount    int
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
		if err := s.ChangeHandler.Handle(s, change); err != nil {
			s.game.ErrorChannel <- err
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
