package session

import (
	"errors"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrUserIDNil   = errors.New("session user id is nil")
	ErrPlayerIDNil = errors.New("player id is nil")
)

type Instance struct {
	UserID    uuid.UUID
	PlayerID  uuid.UUID
	LobbyID   uuid.UUID
	PartySlot int
	DeckIndex int
}

func NewInstance(user uuid.UUID, player uuid.UUID) (*Instance, error) {
	if user == uuid.Nil {
		return nil, ErrUserIDNil
	}
	if player == uuid.Nil {
		return nil, ErrPlayerIDNil
	}
	return &Instance{UserID: user, PlayerID: player}, nil
}

func (x *Instance) CanJoinLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil)
}

func (x *Instance) CanLeaveLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil) == false
}
