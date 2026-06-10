package action

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type HPConsensusChange struct {
	GameID  uuid.UUID
	Enemies []game.EnemyHP
}

func NewHPConsensusChange(gameID uuid.UUID, enemies []game.EnemyHP) HPConsensusChange {
	return HPConsensusChange{GameID: gameID, Enemies: enemies}
}

func (c HPConsensusChange) Identifier() game.ChangeIdentifier {
	return game.ChangeIdentifierHPConsensus
}
