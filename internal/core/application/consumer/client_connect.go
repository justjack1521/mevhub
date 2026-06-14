package consumer

import (
	"mevhub/internal/core/domain/player"
	"time"

	"github.com/justjack1521/mevium/pkg/genproto/protocommon"
	"github.com/justjack1521/mevium/pkg/mevent"
	"github.com/justjack1521/mevrabbit"
	uuid "github.com/satori/go.uuid"
	"github.com/wagslane/go-rabbitmq"
)

type ClientConnectConsumer struct {
	publisher *mevent.Publisher
	*mevrabbit.StandardConsumer
}

func NewClientConnectConsumer(publisher *mevent.Publisher, conn *rabbitmq.Conn) *ClientConnectConsumer {
	var service = &ClientConnectConsumer{
		publisher: publisher,
	}
	consumer, err := mevrabbit.NewStandardConsumer(conn, mevrabbit.ClientConnect, mevrabbit.ClientConnected, mevrabbit.Client, service.Consume)
	if err != nil {
		panic(err)
	}
	service.StandardConsumer = consumer
	return service

}

func (s *ClientConnectConsumer) Consume(ctx *mevrabbit.ConsumerContext) (action rabbitmq.Action, err error) {
	if ctx.UserID() == uuid.Nil || ctx.PlayerID() == uuid.Nil {
		return rabbitmq.NackDiscard, nil
	}
	message, err := protocommon.NewClientConnectedMessage(ctx.Delivery.Body)
	if err != nil {
		return rabbitmq.NackDiscard, nil
	}
	id, err := uuid.FromString(message.SessionId)
	if err != nil {
		return rabbitmq.NackDiscard, nil
	}
	var evt = player.NewConnectedEvent(ctx.Context, id, ctx.UserID(), ctx.PlayerID(), time.Now().UTC())
	s.publisher.Notify(evt)
	return rabbitmq.Ack, nil
}
