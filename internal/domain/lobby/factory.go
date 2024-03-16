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
	FromInvite      bool
	Locked          bool
}

func (f InstanceFactory) Create(id uuid.UUID, party string, options InstanceFactoryOptions) (*Instance, error) {

	instance, err := NewInstance(id, party)
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

	instance.PlayerSlots = make([]PlayerSlot, options.PlayerSlots)

	for index := range instance.PlayerSlots {
		restriction, _ := options.SlotRestrictions[index]
		participant, err := instance.NewPlayerSlot(index, restriction)
		if err != nil {
			return nil, err
		}
		instance.PlayerSlots[index] = participant
	}

	return instance, nil

}

type ParticipantFactory struct {
}

type ParticipantFactoryOptions struct {
	RoleID     uuid.UUID
	SlotIndex  int
	DeckIndex  int
	UseStamina bool
	FromInvite bool
}

func (f ParticipantFactory) Create(client, player uuid.UUID, instance *Instance, options ParticipantFactoryOptions) (*Participant, error) {

	participant, err := instance.NewPlayerParticipant(client, player, options.SlotIndex)
	if err != nil {
		return nil, err
	}

	participant.FromInvite = options.FromInvite

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
