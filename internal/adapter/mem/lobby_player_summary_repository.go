package mem

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type LobbyPlayerSummaryMemoryRepository struct {
	data map[uuid.UUID]lobby.PlayerSummary
}

func NewLobbyPlayerSummaryMemoryRepository() *LobbyPlayerSummaryMemoryRepository {
	return &LobbyPlayerSummaryMemoryRepository{data: make(map[uuid.UUID]lobby.PlayerSummary)}
}

func (r *LobbyPlayerSummaryMemoryRepository) Query(ctx context.Context, id uuid.UUID) (lobby.PlayerSummary, error) {
	summary, exists := r.data[id]
	if exists == false {
		return lobby.PlayerSummary{}, lobby.ErrLobbyPlayerSummaryNotFound(id)
	}
	return summary, nil
}

func (r *LobbyPlayerSummaryMemoryRepository) Create(ctx context.Context, player lobby.PlayerSummary) error {
	r.data[player.Identity.PlayerID] = player
	return nil
}
