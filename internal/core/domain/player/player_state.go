package player

import uuid "github.com/satori/go.uuid"

type Identity struct {
	ClientID uuid.UUID
	PlayerID uuid.UUID
}

type Player struct {
	Identity   Identity
	RoleID     uuid.UUID
	DeckIndex  int
	UseStamina bool
	FromInvite bool
}
