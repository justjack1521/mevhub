package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

type LobbyStampCooldownRepository interface {
	// Claim starts the user's stamp cooldown and reports true, or reports false
	// while a previous one is still running.
	Claim(ctx context.Context, user uuid.UUID, cooldown time.Duration) (bool, error)
}
