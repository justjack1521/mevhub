package app

import (
	"github.com/justjack1521/mevium/pkg/rabbitmv"
	"github.com/wagslane/go-rabbitmq"
)

type ApplicationConsumer interface {
	Consume(ctx *rabbitmv.ConsumerContext) (action rabbitmq.Action, err error)
	Close()
}
