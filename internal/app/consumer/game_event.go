package consumer

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type GameNotification interface {
	mevent.ContextEvent
	GameID() uuid.UUID
	Operation() protomulti.MultiGameNotificationType
	Data() []byte
}

type GameClientNotificationEvent struct {
	ctx       context.Context
	operation protomulti.MultiGameNotificationType
	id        uuid.UUID
	data      []byte
}

func NewGameClientNotificationEvent(ctx context.Context, op protomulti.MultiGameNotificationType, id uuid.UUID, data []byte) GameClientNotificationEvent {
	return GameClientNotificationEvent{ctx: ctx, operation: op, id: id, data: data}
}

func (e GameClientNotificationEvent) Name() string {
	return "game.client.notification"
}

func (e GameClientNotificationEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"operation":  e.operation,
		"game.id":    e.id,
		"length":     len(e.data),
	}
}

func (e GameClientNotificationEvent) Context() context.Context {
	return e.ctx
}

func (e GameClientNotificationEvent) GameID() uuid.UUID {
	return e.id
}

func (e GameClientNotificationEvent) Operation() protomulti.MultiGameNotificationType {
	return e.operation
}

func (e GameClientNotificationEvent) Data() []byte {
	return e.data
}
