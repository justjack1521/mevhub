package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

const hpConsensusTimeout = time.Second * 15

type EnemyTurnState struct {
	QueuedActions map[int][]*game.PlayerActionQueue
	hpReports     map[uuid.UUID][]game.EnemyHP
	startTime     time.Time
	resolved      bool
}

func NewEnemyTurnState(instance *game.LiveGameInstance) *EnemyTurnState {
	var state = &EnemyTurnState{
		QueuedActions: make(map[int][]*game.PlayerActionQueue),
		hpReports:     make(map[uuid.UUID][]game.EnemyHP),
		startTime:     time.Now().UTC(),
	}
	for _, party := range instance.Parties {
		var queue = make([]*game.PlayerActionQueue, party.GetPlayerCount())
		for _, player := range party.Players {
			player.Ready = false
			queue[player.ActionLockIndex] = &game.PlayerActionQueue{
				PlayerID: player.PlayerID,
				Actions:  player.Actions,
			}
		}
		state.QueuedActions[party.PartyIndex] = queue
	}
	return state
}

func (s *EnemyTurnState) submitHP(playerID uuid.UUID, enemies []game.EnemyHP) {
	s.hpReports[playerID] = enemies
}

func (s *EnemyTurnState) allReported(instance *game.LiveGameInstance) bool {
	return len(s.hpReports) >= instance.GetPlayerCount()
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
	return result
}

func (s *EnemyTurnState) finalize(instance *game.LiveGameInstance) {
	if s.resolved {
		return
	}
	s.resolved = true

	if len(s.hpReports) == 0 {
		instance.ActionChannel <- NewStateChangeAction(instance.InstanceID, NewPlayerTurnState(instance))
		return
	}

	consensus := s.calculateConsensus()
	instance.SendChange(NewHPConsensusChange(instance.InstanceID, consensus))

	allDead := true
	for _, e := range consensus {
		if e.HP > 0 {
			allDead = false
			break
		}
	}

	if allDead {
		instance.ActionChannel <- NewStateChangeAction(instance.InstanceID, NewEndGameState(instance))
	} else {
		instance.ActionChannel <- NewStateChangeAction(instance.InstanceID, NewPlayerTurnState(instance))
	}
}

func (s *EnemyTurnState) Update(instance *game.LiveGameInstance, t time.Time) {
	if s.resolved {
		return
	}
	if instance.GetPlayerCount() == 0 {
		return
	}
	if s.allReported(instance) {
		s.finalize(instance)
		return
	}
	if t.Sub(s.startTime) > hpConsensusTimeout {
		s.finalize(instance)
	}
}
