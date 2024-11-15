package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
)

type MatchmakingDispatcher interface {
	DispatchMatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, players []match.PlayerQueueEntry) error
}
