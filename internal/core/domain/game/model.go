package game

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	ErrQuestNotAvailable = errors.New("quest is not available")
	ErrQuestClosingSoon  = errors.New("quest closes too soon to host a new lobby")
)

// LobbyCreateCutoff stops new lobbies being hosted for a quest this close to
// its window closing, so a lobby always has time to fill. Finding and joining
// stay open until the window actually closes. The client mirrors this value.
const LobbyCreateCutoff = 15 * time.Minute

type Mode struct {
	SysID          uuid.UUID
	ModeIdentifier ModeIdentifier
	MaxPlayers     int
	MaxLobbies     int
	FulfillMethod  FulfillMethod
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
	Schedules  []Schedule
}

func (x Quest) Zero() bool {
	return x.SysID == uuid.Nil
}

// AvailableAt fails unless one of the quest's windows is open at now. A quest
// with no windows is never available.
func (x Quest) AvailableAt(now time.Time) error {
	for _, schedule := range x.Schedules {
		if schedule.OpenAt(now) {
			return nil
		}
	}
	return ErrQuestNotAvailable
}

// HostableAt fails unless a window is open at now and stays open for at least
// the cutoff.
func (x Quest) HostableAt(now time.Time, cutoff time.Duration) error {
	if err := x.AvailableAt(now); err != nil {
		return err
	}
	for _, schedule := range x.Schedules {
		if schedule.OpenAt(now) && schedule.OpenAt(now.Add(cutoff)) {
			return nil
		}
	}
	return ErrQuestClosingSoon
}

// Schedule is one availability window for a quest. A zero start or end date
// leaves that side of the window open.
type Schedule struct {
	StartDate time.Time
	EndDate   time.Time
}

func (x Schedule) OpenAt(now time.Time) bool {
	if x.StartDate.IsZero() == false && x.StartDate.After(now) {
		return false
	}
	if x.EndDate.IsZero() == false && x.EndDate.After(now) == false {
		return false
	}
	return true
}

type Category struct {
	SysID uuid.UUID
}

func (x Category) Zero() bool {
	return x == Category{}
}
