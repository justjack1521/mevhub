package lobby

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type SearchEntry struct {
	InstanceID         uuid.UUID
	ModeIdentifier     string
	Level              int
	MinimumPlayerLevel int
	Categories         []uuid.UUID
}

type SearchQuery struct {
	ModeIdentifier     string
	MinimumPlayerLevel int
	Levels             []int
	Categories         []uuid.UUID
}

type SearchResult struct {
	LobbyID uuid.UUID
	Expires time.Duration
}
