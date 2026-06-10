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

func (a *SubmitHPConsensusAction) Perform(instance *game.LiveGameInstance) error {
	state, ok := instance.State.(*EnemyTurnState)
	if !ok {
		return nil
	}
	state.submitHP(a.PlayerID, a.Enemies)
	if state.allReported(instance) {
		state.finalize(instance)
	}
	return nil
}
