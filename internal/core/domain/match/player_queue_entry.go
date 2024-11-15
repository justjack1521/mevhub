package match

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type PlayerQueueEntry struct {
	PlayerID  uuid.UUID
	QuestID   uuid.UUID
	DeckLevel int
	JoinedAt  time.Time
}

func (e PlayerQueueEntry) Zero() bool {
	return e == PlayerQueueEntry{}
}
