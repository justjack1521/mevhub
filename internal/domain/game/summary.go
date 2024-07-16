package game

import uuid "github.com/satori/go.uuid"

type Summary struct {
	SysID        uuid.UUID
	PartyID      string
	Seed         int64
	Participants []PlayerParticipant
}
