package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyCreateCommand struct {
	BasicCommand
	LobbyID   uuid.UUID
	QuestID   uuid.UUID
	PartyID   string
	DeckIndex int
	Comment   string
	Options   CreateLobbyOptions
}

type CreateLobbyOptions struct {
	MinimumPlayerLevel int
	Restrictions       []lobby.PlayerSlotRestriction
}

func (c LobbyCreateCommand) CommandName() string {
	return "lobby.create"
}

func NewLobbyCreateCommand(quest uuid.UUID, deck int, comment string, options CreateLobbyOptions) *LobbyCreateCommand {
	return &LobbyCreateCommand{
		LobbyID:   uuid.NewV4(),
		QuestID:   quest,
		PartyID:   fmt.Sprintf("%08d", rand.Intn(100000000)),
		DeckIndex: deck,
		Comment:   comment,
		Options:   options,
	}
}

type LobbyCreateCommandHandler struct {
	EventPublisher        *mevent.Publisher
	InstanceRepository    port.LobbyInstanceWriteRepository
	QuestRepository       port.QuestRepository
	ParticipantFactory    lobby.ParticipantFactory
	ParticipantRepository port.LobbyParticipantWriteRepository
}

func NewLobbyCreateCommandHandler(publisher *mevent.Publisher, instances port.LobbyInstanceWriteRepository, quests port.QuestRepository, participants port.LobbyParticipantWriteRepository) *LobbyCreateCommandHandler {
	return &LobbyCreateCommandHandler{
		EventPublisher:        publisher,
		InstanceRepository:    instances,
		QuestRepository:       quests,
		ParticipantFactory:    lobby.ParticipantFactory{},
		ParticipantRepository: participants,
	}
}

func (h *LobbyCreateCommandHandler) Handle(ctx Context, cmd *LobbyCreateCommand) error {

	quest, err := h.QuestRepository.QueryByID(cmd.QuestID)
	if err != nil {
		return err
	}

	var factory = lobby.NewInstanceFactory(ctx, ctx.UserID(), ctx.PlayerID())

	var opts = lobby.InstanceFactoryOptions{
		QuestID:            quest.SysID,
		PlayerSlots:        quest.Tier.GameMode.MaxPlayers,
		MinimumPlayerLevel: cmd.Options.MinimumPlayerLevel,
		SlotRestrictions:   make(map[int]lobby.PlayerSlotRestriction),
	}

	for _, value := range cmd.Options.Restrictions {
		opts.SlotRestrictions[value.Index] = value
	}

	instance, err := factory.Create(cmd.LobbyID, cmd.PartyID, opts)
	if err != nil {
		return err
	}

	if err := h.InstanceRepository.Create(ctx, instance); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewInstanceCreatedEvent(ctx, instance.SysID, cmd.QuestID, cmd.PartyID, cmd.Comment, instance.MinimumPlayerLevel))

	for i := 0; i < instance.PlayerSlotCount; i++ {

		var part = lobby.ParticipantJoinOptions{
			RoleID:     uuid.Nil,
			SlotIndex:  i,
			DeckIndex:  0,
			UseStamina: false,
		}

		var user = uuid.Nil
		var player = uuid.Nil

		if i == 0 {
			part.DeckIndex = cmd.DeckIndex
			part.UseStamina = true
			user = ctx.UserID()
			player = ctx.PlayerID()
		}

		participant, err := h.ParticipantFactory.Create(user, player, instance, opts.SlotRestrictions[0], part)
		if err != nil {
			return err
		}

		if err := h.ParticipantRepository.Create(ctx, participant); err != nil {
			return err
		}

		if uuid.Equal(player, uuid.Nil) == false {
			cmd.QueueEvent(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))
		}

	}

	return nil

}
