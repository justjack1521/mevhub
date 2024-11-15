package game

import (
	"mevhub/internal/core/domain/lobby"
	"time"
)

type ModeIdentifier string

const (
	ModeIdentifierNone        = "none"
	ModeIdentifierCoopDefault = "coop_default"
	ModeIdentifierCoopVersus  = "coop_versus"
)

type InstanceOptions struct {
	MinimumPlayerLevel int
	MaxPlayerCount     int
	MaxRunTime         time.Duration
	PlayerTurnDuration time.Duration
	Restrictions       []lobby.PartySlotRestriction
}
