package session

import (
	uuid "github.com/satori/go.uuid"
)

type Instance struct {
	ClientID  uuid.UUID
	PlayerID  uuid.UUID
	LobbyID   uuid.UUID
	PartySlot int
	DeckIndex int
}

func (x *Instance) CanJoinLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil)
}

func (x *Instance) CanLeaveLobby() bool {
	return uuid.Equal(x.LobbyID, uuid.Nil) == false
}
