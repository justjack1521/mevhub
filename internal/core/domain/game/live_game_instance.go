package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
	"sync/atomic"
	"time"
)

const (
	StateTickPeriod = time.Millisecond * 250
	// DisconnectGracePeriod bounds how long a disconnected player's seat is
	// held before the game loop evicts them. The application layer's client
	// reaper uses the same window; the domain deadline is the authoritative
	// backstop that guarantees a turn-boundary stall always lifts.
	DisconnectGracePeriod = time.Minute * 3
)

var (
	ErrFailedPerformAction = func(id uuid.UUID, err error) error {
		return fmt.Errorf("live game %s failed to perform action: %w", id, err)
	}
)

type PartyInstanceOptions struct {
	MaxPlayerCount     int
	PlayerTurnDuration time.Duration
}

// LiveGameInstance is a single-writer actor: Run is the only goroutine that
// may touch Parties, State, or any player field. Everything else communicates
// with the game through ActionChannel and reads through emitted Changes.
type LiveGameInstance struct {
	InstanceID    uuid.UUID
	ActionChannel chan Action
	ChangeChannel chan Change
	ErrorChannel  chan error
	Parties       map[uuid.UUID]*LiveParty
	State         State
	GameDuration  time.Duration
	EndedAt       time.Time
	MaxPartyCount int
	PartyOptions  PartyInstanceOptions
	// LastEnemyHP holds the most recently resolved enemy HP consensus so a
	// reconnecting player can be re-synced with the current enemy state. It is
	// empty until the first enemy turn resolves.
	LastEnemyHP []EnemyHP

	ended    atomic.Bool
	done     chan struct{}
	stopOnce sync.Once
}

func NewLiveGameInstance(source *Instance) *LiveGameInstance {
	var game = &LiveGameInstance{
		InstanceID:    source.SysID,
		ActionChannel: make(chan Action, 64),
		ChangeChannel: make(chan Change, 64),
		ErrorChannel:  make(chan error, 16),
		Parties:       make(map[uuid.UUID]*LiveParty),
		GameDuration:  source.Options.MaxRunTime,
		MaxPartyCount: source.Options.MaxPartyCount,
		PartyOptions: PartyInstanceOptions{
			MaxPlayerCount:     source.Options.MaxPlayerCount,
			PlayerTurnDuration: source.Options.PlayerTurnDuration,
		},
		done: make(chan struct{}),
	}
	return game
}

// End marks the game over. Safe to read from other goroutines via HasEnded.
func (game *LiveGameInstance) End() {
	if game.ended.CompareAndSwap(false, true) {
		game.EndedAt = time.Now().UTC()
	}
}

func (game *LiveGameInstance) HasEnded() bool {
	return game.ended.Load()
}

// Stop signals Run and the server's watcher goroutines to exit. The data
// channels are never closed — senders may still hold references — so shutdown
// is purely done-channel driven and idempotent.
func (game *LiveGameInstance) Stop() {
	game.stopOnce.Do(func() {
		close(game.done)
	})
}

func (game *LiveGameInstance) Done() <-chan struct{} {
	return game.done
}

func (game *LiveGameInstance) GetPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetPlayerCount()
	}
	return total
}

func (game *LiveGameInstance) PartyExists(id uuid.UUID) bool {
	_, exists := game.Parties[id]
	return exists
}

func (game *LiveGameInstance) GetParty(id uuid.UUID) (*LiveParty, error) {
	party, exists := game.Parties[id]
	if exists == false {
		return nil, ErrPlayerNotInGame
	}
	return party, nil
}

func (game *LiveGameInstance) PlayerExists(id uuid.UUID) bool {
	for _, party := range game.Parties {
		if party.PlayerExists(id) {
			return true
		}
	}
	return false
}

func (game *LiveGameInstance) GetPlayer(id uuid.UUID) (*LivePlayer, error) {
	for _, party := range game.Parties {
		player, err := party.GetPlayer(id)
		if err == nil {
			return player, nil
		}
	}
	return nil, ErrPlayerNotInGame
}

func (game *LiveGameInstance) GetPartyForPlayer(id uuid.UUID) (*LiveParty, error) {
	for _, party := range game.Parties {
		if party.PlayerExists(id) {
			return party, nil
		}
	}
	return nil, ErrPlayerNotInGame
}

func (game *LiveGameInstance) GetReadyPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetReadyPlayerCount()
	}
	return total
}

func (game *LiveGameInstance) RemovePlayer(id uuid.UUID) error {
	for _, party := range game.Parties {
		if party.PlayerExists(id) == false {
			continue
		}
		return party.RemovePlayer(id)
	}
	return ErrPlayerNotInGame
}

func (game *LiveGameInstance) HasDisconnectedPlayers() bool {
	for _, party := range game.Parties {
		for _, player := range party.Players {
			if player.Disconnected {
				return true
			}
		}
	}
	return false
}

func (game *LiveGameInstance) GetActionLockedPlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		total += party.GetActionLockedPlayerCount()
	}
	return total
}

// Run is the game's single writer: it applies queued actions and drives state
// ticks on one goroutine, which is what makes the aggregate safe without
// locks. It exits only via Stop.
func (game *LiveGameInstance) Run() {

	ticker := time.NewTicker(StateTickPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-game.done:
			return
		case action := <-game.ActionChannel:
			if action == nil {
				continue
			}
			if err := action.Perform(game); err != nil {
				game.SendError(ErrFailedPerformAction(game.InstanceID, err))
			}
		case t := <-ticker.C:
			if game.HasEnded() {
				continue
			}
			game.State.Update(game, t)
		}
	}

}

// SendChange never blocks past shutdown: a change no consumer will drain is
// dropped once the game is stopped.
func (game *LiveGameInstance) SendChange(change Change) {
	select {
	case game.ChangeChannel <- change:
	case <-game.done:
	}
}

func (game *LiveGameInstance) SendError(err error) {
	select {
	case game.ErrorChannel <- err:
	case <-game.done:
	}
}
