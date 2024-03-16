package query

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type Context struct {
	context.Context
	ClientID uuid.UUID
	PlayerID uuid.UUID
}

func NewContext(context context.Context, client uuid.UUID) *Context {
	return &Context{Context: context, ClientID: client}
}
