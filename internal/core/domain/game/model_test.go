package game

import (
	"errors"
	"testing"
	"time"
)

func TestQuestSchedule(t *testing.T) {

	var now = time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)

	var cases = []struct {
		name      string
		schedules []Schedule
		available error
		hostable  error
	}{
		{"no rows is unavailable", nil, ErrQuestNotAvailable, ErrQuestNotAvailable},
		{"no start or end is always open", []Schedule{{}}, nil, nil},
		{"no start, open until end", []Schedule{{EndDate: now.Add(time.Hour)}}, nil, nil},
		{"no end, open after start", []Schedule{{StartDate: now.Add(-time.Hour)}}, nil, nil},
		{"not yet started", []Schedule{{StartDate: now.Add(time.Hour)}}, ErrQuestNotAvailable, ErrQuestNotAvailable},
		{"already ended", []Schedule{{EndDate: now}}, ErrQuestNotAvailable, ErrQuestNotAvailable},
		{"inside the cutoff can join but not host", []Schedule{{EndDate: now.Add(LobbyCreateCutoff - time.Minute)}}, nil, ErrQuestClosingSoon},
		{"any open window counts", []Schedule{{EndDate: now.Add(-time.Hour)}, {StartDate: now.Add(-time.Minute)}}, nil, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var quest = Quest{Schedules: c.schedules}
			if err := quest.AvailableAt(now); errors.Is(err, c.available) == false {
				t.Errorf("AvailableAt = %v, want %v", err, c.available)
			}
			if err := quest.HostableAt(now, LobbyCreateCutoff); errors.Is(err, c.hostable) == false {
				t.Errorf("HostableAt = %v, want %v", err, c.hostable)
			}
		})
	}

}
