package match

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	ErrNotEnoughLobbiesForGame = func(actual, expected int) error {
		return fmt.Errorf("not enough lobbies for game, expected %d got %d", expected, actual)
	}
)

type LobbyQueueEntryCollection []LobbyQueueEntry

type LobbyQueueEntry struct {
	LobbyID      uuid.UUID
	QuestID      uuid.UUID
	AverageLevel int
	JoinedAt     time.Time
}

func (e LobbyQueueEntry) Zero() bool {
	return e == LobbyQueueEntry{}
}
