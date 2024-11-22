package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

type EndGameState struct {
}

func (s *EndGameState) Update(instance *game.LiveGameInstance, t time.Time) {

}

func NewEndGameState(game *game.LiveGameInstance) *EndGameState {

	game.Ended = true
	game.EndedAt = time.Now().UTC()

	return &EndGameState{}
}
