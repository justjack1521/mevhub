package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
)

type LobbyMatchmakingDispatcher interface {
	Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, lobbies []match.LobbyQueueEntry) error
}

type PlayerMatchmakingDispatcher interface {
	Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, lobby match.LobbyQueueEntry, player match.PlayerQueueEntry) error
}
