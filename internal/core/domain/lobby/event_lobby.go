package lobby

import (
	uuid "github.com/satori/go.uuid"
)

type LobbyEvent interface {
	LobbyID() uuid.UUID
}
