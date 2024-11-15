package consumer

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type LobbyNotification interface {
	mevent.ContextEvent
	LobbyID() uuid.UUID
	Operation() protomulti.MultiLobbyNotificationType
	Data() []byte
}

type LobbyClientNotificationEvent struct {
	ctx       context.Context
	operation protomulti.MultiLobbyNotificationType
	lobby     uuid.UUID
	data      []byte
}

func NewLobbyClientNotificationEvent(ctx context.Context, op protomulti.MultiLobbyNotificationType, id uuid.UUID, data []byte) LobbyClientNotificationEvent {
	return LobbyClientNotificationEvent{ctx: ctx, operation: op, lobby: id, data: data}
}

func (e LobbyClientNotificationEvent) Name() string {
	return "lobby.client.notification"
}

func (e LobbyClientNotificationEvent) ToLogFields() logrus.Fields {
	return logrus.Fields{
		"event.name": e.Name(),
		"operation":  e.operation,
		"lobby.id":   e.lobby,
		"length":     len(e.data),
	}
}

func (e LobbyClientNotificationEvent) Context() context.Context {
	return e.ctx
}

func (e LobbyClientNotificationEvent) LobbyID() uuid.UUID {
	return e.lobby
}

func (e LobbyClientNotificationEvent) Operation() protomulti.MultiLobbyNotificationType {
	return e.operation
}

func (e LobbyClientNotificationEvent) Data() []byte {
	return e.data
}
