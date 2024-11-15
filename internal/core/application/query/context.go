package query

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type Context interface {
	context.Context
	UserID() uuid.UUID
	PlayerID() uuid.UUID
}
