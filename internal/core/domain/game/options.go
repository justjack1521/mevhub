package game

import (
	"mevhub/internal/core/domain/lobby"
	"time"
)

type ModeIdentifier string

const (
	ModeIdentifierNone        = "none"
	ModeIdentifierCoopDefault = "coop_default"
	ModeIdentifierCompSingle  = "comp_solo"
	ModeIdentifierCompDuo     = "comp_duo"
)

type FulfillMethod string

const (
	FulfillMethodNone   = "none"
	FulfillMethodSearch = "search"
	FulfillMethodMatch  = "match"
)

type InstanceOptions struct {
	MinimumPlayerLevel int
	MaxPlayerCount     int
	MaxRunTime         time.Duration
	PlayerTurnDuration time.Duration
	Restrictions       []lobby.PartySlotRestriction
}
