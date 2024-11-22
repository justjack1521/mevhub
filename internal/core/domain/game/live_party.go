package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type LiveParty struct {
	PartyID            uuid.UUID
	PartyIndex         int
	Players            map[uuid.UUID]*LivePlayer
	MaxPlayerCount     int
	PlayerTurnDuration time.Duration
	LastAction         time.Time
}

func (x *LiveParty) GetPlayerCount() int {
	return len(x.Players)
}

func (x *LiveParty) PlayerExists(id uuid.UUID) bool {
	_, exists := x.Players[id]
	return exists
}

func (x *LiveParty) GetPlayer(id uuid.UUID) (*LivePlayer, error) {
	player, exists := x.Players[id]
	if exists == false {
		return nil, ErrPlayerNotInParty
	}
	return player, nil
}

func (x *LiveParty) GetReadyPlayerCount() int {
	var count int
	for _, player := range x.Players {
		if player.Ready {
			count++
		}
	}
	return count
}

func (x *LiveParty) RemovePlayer(id uuid.UUID) error {
	if x.PlayerExists(id) == false {
		return ErrPlayerNotInParty
	}
	delete(x.Players, id)
	return nil
}

func (x *LiveParty) GetActionLockedPlayerCount() int {
	var count int
	for _, player := range x.Players {
		if player.ActionsLocked {
			count++
		}
	}
	return count
}
