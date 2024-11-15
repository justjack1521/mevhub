package lobby

import (
	uuid "github.com/satori/go.uuid"
)

type PartySlotRestriction interface {
	CanJoin(participant Participant) bool
}

type RoleRestriction struct {
	AppliesTo int
	Role      uuid.UUID
}

func (r RoleRestriction) CanJoin(participant Participant) bool {
	return r.Role == participant.Role && participant.PlayerSlot == r.AppliesTo
}

type ReservedRestriction struct {
	AppliesTo int
}

func (r ReservedRestriction) CanJoin(participant Participant) bool {
	return (participant.FromInvite || participant.IsHost()) && participant.PlayerSlot == r.AppliesTo
}

type BotRestriction struct {
	AppliesTo int
}

func (r BotRestriction) CanJoin(participant Participant) bool {
	return participant.PlayerSlot != r.AppliesTo
}

type LockedRestriction struct {
	AppliesTo int
}

func (r LockedRestriction) CanJoin(participant Participant) bool {
	return participant.PlayerSlot != r.AppliesTo
}
