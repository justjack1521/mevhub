package game

import (
	"errors"
)

var (
	ErrPartyIDNil          = errors.New("party id is nil")
	ErrPartyAlreadyInGame  = errors.New("party already added to game")
	ErrPlayerGameFull      = errors.New("live game is full")
	ErrPlayerAlreadyInGame = errors.New("player already added to game")
	ErrPlayerNotInGame     = errors.New("player not in game")
	ErrPlayerNotInParty    = errors.New("player not in party")
	ErrUserIDNil           = errors.New("user id is nil")
	ErrPlayerIDNil         = errors.New("player id is nil")
)

type Action interface {
	Perform(game *LiveGameInstance) error
}
