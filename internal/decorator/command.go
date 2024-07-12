package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Context interface {
	context.Context
}

type Command interface {
	CommandName() string
}

type CommandHandler[CTX Context, C Command] interface {
	Handle(ctx CTX, cmd C) error
}

type QueryHandler[CTX Context, C Command, R any] interface {
	Handle(ctx CTX, cmd C) (R, error)
}

type LoggerCommandDecorator[CTX Context, C Command] struct {
	logger *logrus.Logger
	base   CommandHandler[CTX, C]
}

func NewCommandHandlerWithLogger[CTX Context, C Command](base CommandHandler[CTX, C]) CommandHandler[CTX, C] {
	return LoggerCommandDecorator[CTX, C]{
		base: base,
	}
}

func (h LoggerCommandDecorator[CTX, C]) Handle(ctx CTX, cmd C) error {
	return h.base.Handle(ctx, cmd)
}
