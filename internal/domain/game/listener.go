package game

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type NotificationListener struct {
	GameID   uuid.UUID
	UserID   uuid.UUID
	PlayerID uuid.UUID
}

type NotificationListenerReadRepository interface {
	QueryAllForGame(ctx context.Context, id uuid.UUID) ([]NotificationListener, error)
}

type NotificationListenerWriteRepository interface {
	Create(ctx context.Context, id uuid.UUID, user uuid.UUID) error
}
