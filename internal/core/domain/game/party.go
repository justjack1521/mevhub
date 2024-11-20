package game

import uuid "github.com/satori/go.uuid"

type Party struct {
	SysID     uuid.UUID
	PartyID   string
	Index     int
	PartyName string
}

type PartySummary struct {
	SysID     uuid.UUID
	PartyID   string
	Index     int
	PartyName string
	Players   []Player
}
