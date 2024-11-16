package match

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type LobbyQueueEntry struct {
	LobbyID      uuid.UUID
	QuestID      uuid.UUID
	AverageLevel int
	JoinedAt     time.Time
}

func (e LobbyQueueEntry) Zero() bool {
	return e == LobbyQueueEntry{}
}
