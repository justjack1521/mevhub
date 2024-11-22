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

func (s *EnemyTurnState) Update(game *game.LiveGameInstance, t time.Time) {

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {
		game.ActionChannel <- NewStateChangeAction(game.InstanceID, NewPlayerTurnState(game))
	}

}
