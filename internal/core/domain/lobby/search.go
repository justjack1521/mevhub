package lobby

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"time"
)

type SearchEntry struct {
	InstanceID         uuid.UUID
	ModeIdentifier     game.ModeIdentifier
	Level              int
	MinimumPlayerLevel int
	Categories         []uuid.UUID
}

type SearchQuery struct {
	ModeIdentifier     game.ModeIdentifier
	MinimumPlayerLevel int
	Levels             []int
	Categories         []uuid.UUID
}

type SearchResult struct {
	LobbyID uuid.UUID
	Expires time.Duration
}
