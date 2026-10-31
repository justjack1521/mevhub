package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"sort"
	"time"
)

// enemyTurnTimeout bounds the enemy turn: a client that never reports itself
// ready (crashed mid-animation, wedged) must not stall the turn for everyone.
const enemyTurnTimeout = time.Second * 15

type EnemyTurnState struct {
	QueuedActions map[int][]*game.PlayerActionQueue
	hpReports     map[uuid.UUID][]game.EnemyHP
	// readyWaitFrom is when the ready wait's timeout started counting: the turn's
	// start, or the latest reconnect if anyone came back during it.
	readyWaitFrom time.Time
	resolved      bool
	advanced      bool
}

func NewEnemyTurnState(instance *game.LiveGameInstance) *EnemyTurnState {
	var state = &EnemyTurnState{
		QueuedActions: make(map[int][]*game.PlayerActionQueue),
		hpReports:     make(map[uuid.UUID][]game.EnemyHP),
		readyWaitFrom: time.Now().UTC(),
	}
	for _, party := range instance.Parties {
		var players = make([]*game.LivePlayer, 0, len(party.Players))
		for _, player := range party.Players {
			// Clearing Ready arms the turn's exit condition: each client sets it
			// again once it has played the enemy turn out.
			player.Ready = false
			// A dead player contributes no actions.
			if player.Dead {
				continue
			}
			players = append(players, player)
		}

		// Lock order, packed. Lock indices can have gaps (a player removed after
		// locking), so they are sorted rather than used as positions; clients take
		// the position in this queue as the resolution order. Party slot breaks
		// ties, and anyone somehow unlocked goes last, so map iteration order
		// never decides anything.
		sort.SliceStable(players, func(i, j int) bool {
			if players[i].ActionsLocked != players[j].ActionsLocked {
				return players[i].ActionsLocked
			}
			if players[i].ActionLockIndex != players[j].ActionLockIndex {
				return players[i].ActionLockIndex < players[j].ActionLockIndex
			}
			return players[i].PartySlot < players[j].PartySlot
		})

		var queue = make([]*game.PlayerActionQueue, 0, len(players))
		for _, player := range players {
			// Copy the queue: the live player's Actions slice keeps being
			// mutated after this state captures it.
			var actions = make([]*game.PlayerAction, len(player.Actions))
			copy(actions, player.Actions)
			queue = append(queue, &game.PlayerActionQueue{
				PlayerID: player.PlayerID,
				Actions:  actions,
			})
		}
		state.QueuedActions[party.PartyIndex] = queue
	}
	return state
}

// The HP consensus is disabled: nothing calls submitHP, allReported,
// calculateConsensus or resolve any more, and the enemy turn no longer waits on
// or broadcasts enemy HP. The code below is kept intact so the mechanism can be
// switched back on by restoring the two call sites — SubmitHPConsensusAction.Perform
// and the resolve branch of Update.

func (s *EnemyTurnState) submitHP(playerID uuid.UUID, enemies []game.EnemyHP) {
	s.hpReports[playerID] = enemies
}

func (s *EnemyTurnState) allReported(instance *game.LiveGameInstance) bool {
	// Count only reports from players still in the game: a reporter who has
	// since been removed must not stand in for a present player who has not
	// reported yet. (Their report itself stays in the consensus — it was a
	// valid observation.)
	var reported = 0
	for _, party := range instance.Parties {
		for _, player := range party.Players {
			if _, ok := s.hpReports[player.PlayerID]; ok {
				reported++
			}
		}
	}
	return reported >= instance.GetPlayerCount()
}

func (s *EnemyTurnState) calculateConsensus() []game.EnemyHP {
	// majority vote per enemy — ties resolved toward lower HP (safer)
	votes := make(map[int]map[int]int)
	for _, report := range s.hpReports {
		for _, e := range report {
			if votes[e.EnemyIndex] == nil {
				votes[e.EnemyIndex] = make(map[int]int)
			}
			votes[e.EnemyIndex][e.HP]++
		}
	}
	result := make([]game.EnemyHP, 0, len(votes))
	for index, hpVotes := range votes {
		maxCount := 0
		consensusHP := 0
		for hp, count := range hpVotes {
			if count > maxCount || (count == maxCount && hp < consensusHP) {
				maxCount = count
				consensusHP = hp
			}
		}
		result = append(result, game.EnemyHP{EnemyIndex: index, HP: consensusHP})
	}
	// Deterministic ordering: map iteration would otherwise make the
	// broadcast and the sync snapshot disagree on order for identical data.
	sort.Slice(result, func(i, j int) bool {
		return result[i].EnemyIndex < result[j].EnemyIndex
	})
	return result
}

// resolve runs the enemy turn to completion: it computes and broadcasts the HP
// consensus. This is the turn's work and is never stalled by a disconnect. On a
// wipe it advances straight to the end-game state (no reason to wait for
// reconnects); otherwise it leaves the state resolved-but-not-advanced so the
// hand-off to the next player turn can be held at the boundary.
func (s *EnemyTurnState) resolve(instance *game.LiveGameInstance) {
	if s.resolved {
		return
	}
	s.resolved = true

	if len(s.hpReports) == 0 {
		// No reports this turn (e.g. every player was in their grace window):
		// carry the previous consensus forward, but still broadcast it so
		// clients stay in step at the boundary instead of silently diverging.
		if len(instance.LastEnemyHP) > 0 {
			instance.SendChange(NewHPConsensusChange(instance.InstanceID, instance.LastEnemyHP))
		}
		return
	}

	consensus := s.calculateConsensus()

	// An empty consensus (every report carried no enemies) must not fall
	// through to the wipe check below, where vacuous truth would end the game.
	if len(consensus) == 0 {
		return
	}

	instance.LastEnemyHP = consensus
	instance.SendChange(NewHPConsensusChange(instance.InstanceID, consensus))

	allDead := true
	for _, e := range consensus {
		if e.HP > 0 {
			allDead = false
			break
		}
	}

	if allDead {
		s.advanced = true
		transitionTo(instance, NewEndGameState(instance))
	}
}

// restartReadyWait gives the ready wait a fresh timeout from t. Never moves it
// back, so an out of order reconnect cannot shorten a wait already extended.
func (s *EnemyTurnState) restartReadyWait(t time.Time) {
	if t.After(s.readyWaitFrom) {
		s.readyWaitFrom = t
	}
}

func (s *EnemyTurnState) Update(instance *game.LiveGameInstance, t time.Time) {

	evictExpiredDisconnectedPlayers(instance, t)
	evictExpiredDeadPlayers(instance, t)
	expireLapsedReviveClaims(instance, t)
	restateGamePeriodically(instance, t)

	// All players gone (left or evicted): end the game so the host can
	// reclaim it, rather than idling here forever.
	if instance.GetPlayerCount() == 0 {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	if s.advanced {
		return
	}

	// Everyone dead is a wipe: the game is over. A dead player can only be
	// revived while someone in the party is still standing, so there is nothing
	// to wait for, and ending here gives every client the same outcome even if
	// one of their sims has not reached the wipe yet.
	if instance.GetAlivePlayerCount() == 0 {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	// The enemy turn is a client-side interlude: it ends when every living
	// player reports itself ready again (NewEnemyTurnState cleared the flag on
	// entry), or when the timeout lifts a turn a silent client would stall.
	// Dead players are exempt — the turn must not wait on a corpse. The timeout
	// counts from the latest reconnect if there was one: a player who was away
	// past it is back still playing the round out, and would otherwise be handed
	// a player turn the moment they returned.
	if instance.AllAlivePlayersReady() == false && t.Sub(s.readyWaitFrom) <= enemyTurnTimeout {
		return
	}

	// Out of time ends the game at the boundary, never mid-turn: every client
	// has played the dispatched queue and the enemy turn out by now, so all of
	// them leave from the same point. Not held for a disconnected player, as
	// there is no next turn for them to rejoin.
	if instance.TimeExpired(t) {
		transitionTo(instance, NewEndGameState(instance))
		return
	}

	// We are at the enemy -> player boundary. Hold the hand-off to the next
	// player turn while any player is within their reconnect grace period. The
	// stall lifts on reconnect or timeout removal.
	if instance.HasDisconnectedPlayers() {
		return
	}
	s.advanced = true
	transitionTo(instance, NewPlayerTurnState(instance))
}
