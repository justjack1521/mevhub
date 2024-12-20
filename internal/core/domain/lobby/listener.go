package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type ListenerLevel int

const (
	ListenerLevelHost = iota
	ListenerLevelPlayer
	ListenerLevelViewer
)

type NotificationListener struct {
	LobbyID  uuid.UUID
	UserID   uuid.UUID
	PlayerID uuid.UUID
	Level    ListenerLevel
}

type NotificationListenerReadRepository interface {
	QueryAllForLobby(ctx context.Context, id uuid.UUID) ([]NotificationListener, error)
}

type NotificationListenerWriteRepository interface {
	CreateListener(ctx context.Context, lobby uuid.UUID, user uuid.UUID) error
	DeleteListener(ctx context.Context, id uuid.UUID, user uuid.UUID) error
	DeleteAll(ctx context.Context, lobby uuid.UUID) error
}

type NotificationListenerRepository interface {
	NotificationListenerReadRepository
	NotificationListenerWriteRepository
}

func NewNotificationListener(instance, user, player uuid.UUID, level ListenerLevel) *NotificationListener {
	return &NotificationListener{LobbyID: instance, UserID: user, PlayerID: player, Level: level}
}
