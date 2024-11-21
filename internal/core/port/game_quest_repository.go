package port

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

var ErrFailedGetQuestByID = func(id uuid.UUID, err error) error {
	return fmt.Errorf("failed to get quest by id %s: %w", id, err)
}

type QuestRepository interface {
	QueryByID(id uuid.UUID) (game.Quest, error)
}
