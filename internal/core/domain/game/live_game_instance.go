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
	// StateRestatePeriod bounds how long a client can stay wrong. The game
	// restates itself in full at least this often, so a notification lost
	// anywhere between here and the client's socket repairs itself without the
	// client having to notice, ask, or be correct about anything.
	StateRestatePeriod = time.Second * 5
	// DisconnectGracePeriod bounds how long a disconnected player's seat is
	// held before the game loop evicts them, and so how long a turn boundary
	// stalls for a player who is never coming back. The application layer's
	// client reaper uses the same window; the domain deadline is the
	// authoritative backstop that guarantees the stall always lifts.
	DisconnectGracePeriod = time.Second * 30
	// ReviveClaimDuration is how long a granted revive claim holds before it
	// decays. It must comfortably cover the claimant's round trip to consume
	// the revive item (BattleRevive on the game service) plus the PlayerRevive
	// call back here, but stay short enough that an abandoned claim only
	// briefly blocks other would-be revivers.
	ReviveClaimDuration = time.Second * 10
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
	// DeadPlayerKickDuration is how long a player may stay dead with no revive
	// before the death sweep removes them from the game. Zero disables the
	// kick, which is the current configuration everywhere: nothing populates
	// the option yet, so dead players are never kicked.
	DeadPlayerKickDuration time.Duration
	// LastEnemyHP holds the most recently resolved enemy HP consensus so a
	// reconnecting player can be re-synced with the current enemy state. The
	// consensus is disabled, so nothing writes this and the restatement
	// carries no enemies — clients source enemy HP themselves.
	LastEnemyHP []EnemyHP
	// LastSyncAt is when the game last restated itself in full. Read and
	// written only inside the single-writer loop, so it needs no
	// synchronisation of its own.
	LastSyncAt time.Time

	ended    atomic.Bool
	done     chan struct{}
	stopOnce sync.Once
}

func NewLiveGameInstance(source *Instance) *LiveGameInstance {
	var game = &LiveGameInstance{
		InstanceID:             source.SysID,
		ActionChannel:          make(chan Action, 64),
		ChangeChannel:          make(chan Change, 64),
		ErrorChannel:           make(chan error, 16),
		Parties:                make(map[uuid.UUID]*LiveParty),
		GameDuration:           source.Options.MaxRunTime,
		MaxPartyCount:          source.Options.MaxPartyCount,
		DeadPlayerKickDuration: source.Options.DeadPlayerKickDuration,
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

func (game *LiveGameInstance) GetAlivePlayerCount() int {
	var total = 0
	for _, party := range game.Parties {
		for _, player := range party.Players {
			if player.Dead == false {
				total++
			}
		}
	}
	return total
}

// AllAlivePlayersLocked reports whether every living player has locked in
// their action queue. Dead players are exempt: they cannot act, so a turn
// boundary must never wait on them. Vacuously true with no living players —
// callers guard on GetAlivePlayerCount first.
func (game *LiveGameInstance) AllAlivePlayersLocked() bool {
	for _, party := range game.Parties {
		for _, player := range party.Players {
			if player.Dead == false && player.ActionsLocked == false {
				return false
			}
		}
	}
	return true
}

// AllAlivePlayersReady reports whether every living player has flagged itself
// ready. Dead players are exempt for the same reason as AllAlivePlayersLocked,
// and it is likewise vacuously true with no living players.
func (game *LiveGameInstance) AllAlivePlayersReady() bool {
	for _, party := range game.Parties {
		for _, player := range party.Players {
			if player.Dead == false && player.Ready == false {
				return false
			}
		}
	}
	return true
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
//
// It stays blocking by design. A dropped change is a notification no client
// ever sees, and while the next restatement would eventually paper over the
// resulting divergence, dropping under load is how a brief publisher stall
// turns into every client being wrong at once.
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
