package app

import (
	"github.com/justjack1521/mevrabbit"
	"github.com/wagslane/go-rabbitmq"
)

type ApplicationConsumer interface {
	Consume(ctx *mevrabbit.ConsumerContext) (action rabbitmq.Action, err error)
	Close()
}
