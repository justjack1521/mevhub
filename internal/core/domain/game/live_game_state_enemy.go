package game

import "time"

type EnemyTurnState struct {
	QueuedActions map[int][]*PlayerActionQueue
}

func NewEnemyTurnState(game *LiveGameInstance) *EnemyTurnState {
	var state = &EnemyTurnState{
		QueuedActions: make(map[int][]*PlayerActionQueue),
	}
	for _, party := range game.Parties {
		var queue = make([]*PlayerActionQueue, game.GetPlayerCount())
		for _, player := range party.Players {
			player.Ready = false
			queue[player.ActionLockIndex] = &PlayerActionQueue{
				PlayerID: player.PlayerID,
				Actions:  player.Actions,
			}
		}
		state.QueuedActions[party.PartyIndex] = queue
	}
	return state
}

func (s *EnemyTurnState) Update(game *LiveGameInstance, t time.Time) {

	if game.GetPlayerCount() == 0 {
		return
	}

	if game.GetReadyPlayerCount() == game.GetPlayerCount() {
		game.ActionChannel <- NewStateChangeAction(game.InstanceID, NewPlayerTurnState(game))
	}

}
