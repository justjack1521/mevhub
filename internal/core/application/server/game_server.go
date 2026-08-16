package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// ClientTimeoutPeriod mirrors the domain's grace window: the reaper here
	// drives session cleanup, while the live game evicts the player itself via
	// game.DisconnectGracePeriod.
	ClientTimeoutPeriod = game.DisconnectGracePeriod
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
	errorCount    atomic.Int64
	// notifications holds every broadcast this game has published, so a client
	// that missed some can ask for them back. It has its own lock; see
	// notificationLog.
	notifications notificationLog
}

// ClaimExpiredClients returns the channels whose grace period has elapsed and
// marks them so each is only ever returned once. It is safe to call from the
// host goroutine: access to the clients map and the DisconnectedAt/timedOut
// fields is synchronised against the change-handling goroutine.
func (s *GameServer) ClaimExpiredClients(timeout time.Duration) []*PlayerChannel {
	s.mu.Lock()
	defer s.mu.Unlock()
	var expired []*PlayerChannel
	for _, ch := range s.clients {
		if ch.timedOut || ch.DisconnectedAt == nil || time.Since(*ch.DisconnectedAt) < timeout {
			continue
		}
		ch.timedOut = true
		expired = append(expired, ch)
	}
	return expired
}

func (s *GameServer) Start() {
	go s.WatchChanges()
	go s.WatchErrors()
	go s.game.Run()
}

// Stop signals the game loop and both watcher goroutines to exit.
func (s *GameServer) Stop() {
	s.game.Stop()
}

func (s *GameServer) WatchErrors() {
	for {
		select {
		case <-s.game.Done():
			return
		case err := <-s.game.ErrorChannel:
			s.ErrorHandler.Handle(s, err)
		}
	}
}

func (s *GameServer) WatchChanges() {
	for {
		select {
		case <-s.game.Done():
			return
		case change := <-s.game.ChangeChannel:
			if err := s.ChangeHandler.Handle(s, change); err != nil {
				s.game.SendError(err)
			}
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
