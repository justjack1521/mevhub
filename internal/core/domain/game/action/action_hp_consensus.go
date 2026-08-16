package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type SubmitHPConsensusAction struct {
	GameID   uuid.UUID
	PlayerID uuid.UUID
	Enemies  []game.EnemyHP
}

func NewSubmitHPConsensusAction(gameID uuid.UUID, playerID uuid.UUID, enemies []game.EnemyHP) *SubmitHPConsensusAction {
	return &SubmitHPConsensusAction{GameID: gameID, PlayerID: playerID, Enemies: enemies}
}

// Perform is a no-op: the HP consensus is disabled. Reports are still accepted
// over the wire so existing clients do not error, they simply carry no effect —
// the enemy turn now advances on player ready reports alone. To re-enable,
// restore the body below.
func (a *SubmitHPConsensusAction) Perform(instance *game.LiveGameInstance) error {
	return nil
	/*
		state, ok := instance.State.(*EnemyTurnState)
		if !ok {
			return nil
		}
		state.submitHP(a.PlayerID, a.Enemies)
		if state.allReported(instance) {
			state.resolve(instance)
		}
		return nil
	*/
}
