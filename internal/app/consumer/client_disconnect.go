package consumer

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/justjack1521/mevrabbit"
	"github.com/wagslane/go-rabbitmq"
	"mevhub/internal/domain/player"
	"time"
)

type ClientDisconnectConsumer struct {
	publisher *mevent.Publisher
	*mevrabbit.StandardConsumer
}

func NewClientDisconnectConsumer(publisher *mevent.Publisher, conn *rabbitmq.Conn) *ClientDisconnectConsumer {
	var service = &ClientDisconnectConsumer{
		publisher: publisher,
	}
	consumer, err := mevrabbit.NewStandardConsumer(conn, mevrabbit.ClientUpdate, mevrabbit.ClientDisconnected, mevrabbit.Client, service.Consume)
	if err != nil {
		panic(err)
	}
	service.StandardConsumer = consumer
	return service

}

func (s *ClientDisconnectConsumer) Consume(ctx *mevrabbit.ConsumerContext) (action rabbitmq.Action, err error) {
	var evt = player.NewDisconnectedEvent(ctx.Context, ctx.UserID(), ctx.PlayerID(), time.Now().UTC())
	s.publisher.Notify(evt)
	return rabbitmq.Ack, nil
}
