package server

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type PlayerChannel struct {
	UserID      uuid.UUID
	PlayerID    uuid.UUID
	LastMessage time.Time
}

func NewPlayerChannel(user uuid.UUID, player uuid.UUID) *PlayerChannel {
	return &PlayerChannel{UserID: user, PlayerID: player}
}
