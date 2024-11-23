package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

type EnemyTurnState struct {
	QueuedActions map[int][]*game.PlayerActionQueue
}

func NewEnemyTurnState(instance *game.LiveGameInstance) *EnemyTurnState {
	var state = &EnemyTurnState{
		QueuedActions: make(map[int][]*game.PlayerActionQueue),
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

func (s *EnemyTurnState) Update(instance *game.LiveGameInstance, t time.Time) {

	if instance.GetPlayerCount() == 0 {
		return
	}

	if instance.GetReadyPlayerCount() == instance.GetPlayerCount() {
		instance.ActionChannel <- NewStateChangeAction(instance.InstanceID, NewPlayerTurnState(instance))
	}

}
