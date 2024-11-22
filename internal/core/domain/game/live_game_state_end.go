package game

import "time"

type EndGameState struct {
}

func (s *EndGameState) Update(game *LiveGameInstance, t time.Time) {

}

func NewEndGameState(game *LiveGameInstance) *EndGameState {

	game.Ended = true
	game.EndedAt = time.Now().UTC()

	return &EndGameState{}
}
