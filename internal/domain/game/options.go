package game

import (
	"mevhub/internal/domain/lobby"
	"time"
)

type ModeIdentifier string

const (
	ModeIdentifierNone        = "none"
	ModeIdentifierCoopDefault = "coop_default"
)

type InstanceOptions struct {
	MinimumPlayerLevel int
	MaxRunTime         time.Duration
	PlayerTurnDuration time.Duration
	Restrictions       []lobby.PartySlotRestriction
}
