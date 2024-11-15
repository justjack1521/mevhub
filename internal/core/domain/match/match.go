package match

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Instance struct {
	SysID        uuid.UUID
	QuestID      uuid.UUID
	Started      bool
	RegisteredAt time.Time
}

type Participant struct {
	UserID     uuid.UUID
	PlayerID   uuid.UUID
	LobbyID    uuid.UUID
	TeamIndex  int
	PlayerSlot int
	DeckIndex  int
	Ready      bool
}
