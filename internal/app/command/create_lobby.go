package command

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
)

type CreateLobbyCommand struct {
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

func (c CreateLobbyCommand) CommandName() string {
	return "create.lobby"
}

func NewCreateLobbyCommand(quest uuid.UUID, deck int, comment string, options CreateLobbyOptions) CreateLobbyCommand {
	return CreateLobbyCommand{
		LobbyID:   uuid.NewV4(),
		QuestID:   quest,
		PartyID:   fmt.Sprintf("%08d", rand.Intn(100000000)),
		DeckIndex: deck,
		Comment:   comment,
		Options:   options,
	}
}

type CreateLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	InstanceRepository    lobby.InstanceWriteRepository
	QuestRepository       game.QuestRepository
	ParticipantFactory    lobby.ParticipantFactory
	ParticipantRepository lobby.ParticipantWriteRepository
}

func NewCreateLobbyCommandHandler(publisher *mevent.Publisher, instances lobby.InstanceWriteRepository, quests game.QuestRepository, participants lobby.ParticipantWriteRepository) *CreateLobbyCommandHandler {
	return &CreateLobbyCommandHandler{
		EventPublisher:        publisher,
		InstanceRepository:    instances,
		QuestRepository:       quests,
		ParticipantFactory:    lobby.ParticipantFactory{},
		ParticipantRepository: participants,
	}
}

func (h *CreateLobbyCommandHandler) Handle(ctx Context, cmd *CreateLobbyCommand) error {

	quest, err := h.QuestRepository.QueryByID(cmd.QuestID)
	if err != nil {
		return err
	}

	var factory = lobby.NewInstanceFactory(ctx, ctx.UserID())

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

	participant, err := h.ParticipantFactory.Create(ctx.UserID(), ctx.PlayerID(), instance, lobby.ParticipantFactoryOptions{
		RoleID:     uuid.Nil,
		SlotIndex:  0,
		DeckIndex:  cmd.DeckIndex,
		UseStamina: true,
	})
	if err != nil {
		return err
	}

	if err := h.ParticipantRepository.Create(ctx, participant); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantCreatedEvent(ctx, participant.ClientID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	return nil

}
