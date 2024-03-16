package mem

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/game"
)

type GameInstanceMemoryRepository struct {
	data map[uuid.UUID]*game.Instance
}

func NewGameInstanceMemoryRepository() *GameInstanceMemoryRepository {
	return &GameInstanceMemoryRepository{data: make(map[uuid.UUID]*game.Instance)}
}

func (r *GameInstanceMemoryRepository) QueryByID(ctx context.Context, id uuid.UUID) (*game.Instance, error) {
	result, exists := r.data[id]
	if exists == false {
		return nil, game.ErrGameInstanceNotFound(game.ErrGameInstanceNotFoundByID(id))
	}
	return result, nil
}

func (r *GameInstanceMemoryRepository) Create(ctx context.Context, instance *game.Instance) error {
	if instance == nil || instance.SysID == uuid.Nil {
		return game.ErrGameInstanceNil
	}
	r.data[instance.SysID] = instance
	return nil
}
