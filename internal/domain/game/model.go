package game

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Mode struct {
	SysID             uuid.UUID
	OptionsIdentifier ModeIdentifier
	MaxPlayers        int
}

func (x Mode) Zero() bool {
	return x == Mode{}
}

type Tier struct {
	SysID              uuid.UUID
	GameMode           Mode
	StarLevel          int
	StaminaCost        int
	PlayerTurnDuration time.Duration
	TimeLimit          time.Duration
	SeedMultiplier     int
}

func (x Tier) Zero() bool {
	return x == Tier{}
}

type Quest struct {
	SysID      uuid.UUID
	Tier       Tier
	Categories []Category
}

func (x Quest) Zero() bool {
	return x.SysID == uuid.Nil
}

type Category struct {
	SysID uuid.UUID
}

func (x Category) Zero() bool {
	return x == Category{}
}

type QuestRepository interface {
	QueryByID(id uuid.UUID) (Quest, error)
}
