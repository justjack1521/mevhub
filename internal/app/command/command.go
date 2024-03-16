package command

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/session"
)

type Context struct {
	context.Context
	ClientID uuid.UUID
	Session  *session.Instance
}

func NewContext(ctx context.Context, client uuid.UUID) *Context {
	return &Context{Context: ctx, ClientID: client}
}
