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

// GameStateSyncChange is a full snapshot of the live game state targeted at a
// single player (TargetPlayerID), used to re-sync a player who has reconnected
// after their grace period. It is delivered only to the target, not broadcast.
type GameStateSyncChange struct {
	InstanceID      uuid.UUID
	TargetPlayerID  uuid.UUID
	Phase           GameSyncPhase
	TurnRemainingMs int64
	Parties         []GameSyncParty
	Enemies         []game.EnemyHP
}

// NewGameStateSyncChange builds a snapshot from the live instance. It must be
// called from inside the game's single-writer loop (an Action or State), which
// is the only context where the aggregate may be read.
func NewGameStateSyncChange(instance *game.LiveGameInstance, target uuid.UUID) *GameStateSyncChange {

	// Copy: the snapshot crosses the change channel to the publisher
	// goroutine while the loop may resolve a new consensus into LastEnemyHP.
	var enemies = make([]game.EnemyHP, len(instance.LastEnemyHP))
	copy(enemies, instance.LastEnemyHP)

	var change = &GameStateSyncChange{
		InstanceID:     instance.InstanceID,
		TargetPlayerID: target,
		Enemies:        enemies,
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
		} else {
			// TurnDuration == 0 means the turn never expires; 0ms would read
			// as "already over", so signal "no timer" instead.
			change.TurnRemainingMs = -1
		}
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
