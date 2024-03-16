package mem

import (
	"context"
	"mevhub/internal/domain/game"
	"time"
)

type GameInstanceOptionsMemoryRepository struct {
	data map[game.ModeIdentifier]game.InstanceOptions
}

func NewGameInstanceOptionsMemoryRepository() *GameInstanceOptionsMemoryRepository {
	return &GameInstanceOptionsMemoryRepository{data: map[game.ModeIdentifier]game.InstanceOptions{
		game.ModeIdentifierCoopDefault: {
			MinimumPlayerLevel: 120,
			MaxRunTime:         time.Minute * 30,
		},
	}}
}

func (r *GameInstanceOptionsMemoryRepository) QueryByIdentifier(ctx context.Context, identifier game.ModeIdentifier) (game.InstanceOptions, error) {
	result, exists := r.data[identifier]
	if exists == false {
		return game.InstanceOptions{}, game.ErrInstanceOptionsNotFound(game.ErrInstanceOptionsNotFoundWithIdentifier(identifier))
	}
	return result, nil
}
