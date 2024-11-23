package action

import (
	"mevhub/internal/core/domain/game"
	"time"
)

type EndGameState struct {
}

func (s *EndGameState) Update(instance *game.LiveGameInstance, t time.Time) {

}

func NewEndGameState(instance *game.LiveGameInstance) *EndGameState {

	instance.Ended = true
	instance.EndedAt = time.Now().UTC()

	return &EndGameState{}
}
