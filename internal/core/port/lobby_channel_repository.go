package port

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type LobbyNotificationChannelOpener interface {
	Open(ctx context.Context, id uuid.UUID)
}
