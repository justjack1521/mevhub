package server

import (
	"context"
	"errors"
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
	// sequence numbers this game's broadcasts so a client can tell that it lost
	// one. Nothing resends by sequence — the game restates itself instead — so
	// this exists purely so loss is measurable rather than silent.
	sequence atomic.Uint64
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

// ErrGameServerSeedOverflow is a seed that would block: its change could not be
// queued because the change channel is already full and nothing drains it yet.
var ErrGameServerSeedOverflow = errors.New("game server seed would overflow the change channel before start")

// Seed performs an action on the game before Start, while nothing else can
// touch the aggregate. It is how a new game is built populated: its parties and
// players exist before the loop runs, so no client action can reach the game
// ahead of the player it is for. The changes it emits wait on the change channel
// and are handled in order once Start runs - the client registry first.
//
// Only for building a new game. Nothing drains the change channel yet, so a seed
// is refused once the channel is full rather than blocking the caller forever.
func (s *GameServer) Seed(a game.Action) error {
	if len(s.game.ChangeChannel) == cap(s.game.ChangeChannel) {
		return ErrGameServerSeedOverflow
	}
	return a.Perform(s.game)
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
