package decorator

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/sirupsen/logrus"
)

type Context interface {
	context.Context
}

type Query interface {
	CommandName() string
}

type Command interface {
	QueueEvent(evt mevent.Event)
	GetQueuedEvents() []mevent.Event
}

type CommandHandler[CTX Context, C Command] interface {
	Handle(ctx CTX, cmd C) error
}

type QueryHandler[CTX Context, C Query, R any] interface {
	Handle(ctx CTX, cmd C) (R, error)
}

func NewStandardCommandDecorator[CTX Context, C Command](publisher *mevent.Publisher, handler CommandHandler[CTX, C]) CommandHandler[CTX, C] {
	return NewCommandHandlerWithEventPublisher[CTX, C](publisher, handler)
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

type EventPublisherCommandDecorator[CTX Context, C Command] struct {
	publisher *mevent.Publisher
	base      CommandHandler[CTX, C]
}

func NewCommandHandlerWithEventPublisher[CTX Context, C Command](publisher *mevent.Publisher, handler CommandHandler[CTX, C]) CommandHandler[CTX, C] {
	return EventPublisherCommandDecorator[CTX, C]{
		publisher: publisher,
		base:      handler,
	}
}

func (h EventPublisherCommandDecorator[CTX, C]) Handle(ctx CTX, cmd C) (failure error) {
	defer func() {
		if failure == nil {
			for _, value := range cmd.GetQueuedEvents() {
				h.publisher.Notify(value)
			}
		}
	}()
	failure = h.base.Handle(ctx, cmd)
	return
}
