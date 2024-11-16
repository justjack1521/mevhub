package port

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
)

type QuestRepository interface {
	QueryByID(id uuid.UUID) (game.Quest, error)
}
