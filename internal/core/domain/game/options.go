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
	MaxPartyCount      int
	MaxPlayerCount     int
	MaxRunTime         time.Duration
	PlayerTurnDuration time.Duration
	// DeadPlayerKickDuration is how long a player may stay dead with no revive
	// before they are kicked from the game. Deliberately left zero for now (no
	// kick): the quest tier data does not carry it yet, so the factory never
	// populates it. The proto and read-model round trip carry it fine.
	DeadPlayerKickDuration time.Duration
	Restrictions           []lobby.PartySlotRestriction
}
