package server

import uuid "github.com/satori/go.uuid"

type PlayerRegisterNotification struct {
	InstanceID uuid.UUID
	Player     *PlayerChannel
}

type PlayerReadyNotification struct {
	InstanceID uuid.UUID
	PlayerID   uuid.UUID
}
