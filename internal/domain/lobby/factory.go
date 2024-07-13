package lobby

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type InstanceFactory struct {
	ctx    context.Context
	client uuid.UUID
	player uuid.UUID
}

func NewInstanceFactory(ctx context.Context, client uuid.UUID) InstanceFactory {
	return InstanceFactory{ctx: ctx, client: client}
}

type InstanceFactoryOptions struct {
	QuestID            uuid.UUID
	PlayerSlots        int
	MinimumPlayerLevel int
	SlotRestrictions   map[int]PlayerSlotRestriction
}

type PlayerSlotRestriction struct {
	Index           int
	RoleRestriction uuid.UUID
	InviteOnly      bool
	Locked          bool
	BotControl      bool
}

func (f InstanceFactory) Create(id uuid.UUID, party string, options InstanceFactoryOptions) (*Instance, error) {

	instance, err := NewInstance(id, party, options)
	if err != nil {
		return nil, err
	}

	instance.HostID = f.client

	if err := instance.SetQuestID(options.QuestID); err != nil {
		return nil, err
	}

	if err := instance.SetMinPlayerLevel(options.MinimumPlayerLevel); err != nil {
		return nil, err
	}

	return instance, nil

}

type ParticipantFactory struct {
}

type ParticipantJoinOptions struct {
	RoleID     uuid.UUID
	SlotIndex  int
	DeckIndex  int
	UseStamina bool
	FromInvite bool
}

func (f ParticipantFactory) Create(client, player uuid.UUID, instance *Instance, restriction PlayerSlotRestriction, options ParticipantJoinOptions) (*Participant, error) {

	participant, err := instance.NewPlayerParticipant(client, player, restriction, options)
	if err != nil {
		return nil, err
	}

	if err := participant.SetRole(player, options.RoleID); err != nil {
		return nil, err
	}

	if err := participant.SetDeckIndex(player, options.DeckIndex); err != nil {
		return nil, err
	}

	if err := participant.SetUseStamina(player, options.UseStamina); err != nil {
		return nil, err
	}

	if err := participant.SetReady(player, options.SlotIndex == 0); err != nil {
		return nil, err
	}

	if err := instance.CanAddParticipant(participant); err != nil {
		return nil, err
	}

	return participant, nil
}
