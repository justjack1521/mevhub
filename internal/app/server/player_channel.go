package server

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type PlayerChannel struct {
	UserID      uuid.UUID
	PlayerID    uuid.UUID
	Ready       bool
	participant *game.PlayerParticipant
}

func NewPlayerChannel(user uuid.UUID, player uuid.UUID, participant *game.PlayerParticipant) *PlayerChannel {
	return &PlayerChannel{UserID: user, PlayerID: player, participant: participant}
}
