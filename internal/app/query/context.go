package query

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/session"
)

type Context interface {
	context.Context
	UserID() uuid.UUID
	PlayerID() uuid.UUID
	Session() *session.Instance
	SetSession(instance *session.Instance)
}
