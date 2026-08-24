package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

// GameSyncPhase mirrors protomulti.GameSyncPhase; the integer values are kept in
// sync so the marshaller can cast directly.
type GameSyncPhase int

const (
	GameSyncPhaseUnknown GameSyncPhase = iota
	GameSyncPhasePlayerTurn
	GameSyncPhaseEnemyTurn
	// GameSyncPhasePending and GameSyncPhaseEndGame are ahead of the pinned
	// protomulti.GameSyncPhase enum (proto3 open enums carry them fine); the
	// proto definition needs the matching values added before clients can
	// name them. Without these a player syncing during pending or after the
	// game ended was indistinguishable from a server bug (phase 0, 0ms).
	GameSyncPhasePending
	GameSyncPhaseEndGame
)

type GameSyncPlayer struct {
	PlayerID     uuid.UUID
	PlayerIndex  int
	Ready        bool
	Locked       bool
	LockIndex    int
	Disconnected bool
	Actions      []*game.PlayerAction
}

type GameSyncParty struct {
	PartyIndex int
	PartyID    uuid.UUID
	Players    []GameSyncPlayer
}

// GameStateSyncChange restates the whole of the live game state and is
// broadcast to every player. It carries no history: the game holds a roster,
// each player's queued actions and flags, the phase and the turn timer, and
// that set fully describes every notification a client could have missed. A
// client that lost one therefore does not need it resent — it needs this.
type GameStateSyncChange struct {
	InstanceID      uuid.UUID
	Phase           GameSyncPhase
	TurnRemainingMs int64
	Parties         []GameSyncParty
	Enemies         []game.EnemyHP
}

// NewGameStateSyncChange builds the restatement from the live instance. It must
// be called from inside the game's single-writer loop (an Action or State),
// which is the only context where the aggregate may be read.
func NewGameStateSyncChange(instance *game.LiveGameInstance) *GameStateSyncChange {

	// Copy: the restatement crosses the change channel to the publisher
	// goroutine while the loop may resolve a new consensus into LastEnemyHP.
	var enemies = make([]game.EnemyHP, len(instance.LastEnemyHP))
	copy(enemies, instance.LastEnemyHP)

	var change = &GameStateSyncChange{
		InstanceID: instance.InstanceID,
		Enemies:    enemies,
		// -1 is the established "no timer" signal, and it is the correct
		// default for every phase that is not a timed player turn. Leaving the
		// zero value here would report 0ms, which reads as "already over".
		TurnRemainingMs: -1,
	}

	switch state := instance.State.(type) {
	case *PlayerTurnState:
		change.Phase = GameSyncPhasePlayerTurn
		if state.TurnDuration > 0 {
			var remaining = state.TurnDuration - time.Since(state.StartTime)
			if remaining < 0 {
				remaining = 0
			}
			change.TurnRemainingMs = remaining.Milliseconds()
		}
		// TurnDuration == 0 means the turn never expires, so the -1 default
		// stands.
	case *EnemyTurnState:
		change.Phase = GameSyncPhaseEnemyTurn
	case *PendingState:
		change.Phase = GameSyncPhasePending
	case *EndGameState:
		change.Phase = GameSyncPhaseEndGame
	}

	for _, party := range instance.Parties {
		var summary = GameSyncParty{PartyIndex: party.PartyIndex, PartyID: party.PartyID}
		for _, player := range party.Players {
			// Copy the action slice so the snapshot is not mutated by later
			// enqueue/dequeue once it has crossed the change channel.
			var actions = make([]*game.PlayerAction, len(player.Actions))
			copy(actions, player.Actions)
			summary.Players = append(summary.Players, GameSyncPlayer{
				PlayerID:     player.PlayerID,
				PlayerIndex:  player.PartySlot,
				Ready:        player.Ready,
				Locked:       player.ActionsLocked,
				LockIndex:    player.ActionLockIndex,
				Disconnected: player.Disconnected,
				Actions:      actions,
			})
		}
		change.Parties = append(change.Parties, summary)
	}

	return change
}

func (c GameStateSyncChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierGameSync
}
