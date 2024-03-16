package command

import (
	"errors"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/domain/lobby"
)

type JoinLobbyCommand struct {
	LobbyID    uuid.UUID
	DeckIndex  int
	SlotIndex  int
	UseStamina bool
	FromInvite bool
}

func (c JoinLobbyCommand) CommandName() string {
	return "join.lobby"
}

func NewJoinLobbyCommand(id uuid.UUID, deck, slot int, stamina, invite bool) JoinLobbyCommand {
	return JoinLobbyCommand{
		LobbyID:    id,
		DeckIndex:  deck,
		SlotIndex:  slot,
		UseStamina: stamina,
		FromInvite: invite,
	}
}

type JoinLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	InstanceRepository    lobby.InstanceRepository
	ParticipantFactory    lobby.ParticipantFactory
	ParticipantRepository lobby.ParticipantRepository
}

func NewJoinLobbyCommandHandler(publishes *mevent.Publisher, instances lobby.InstanceRepository, participants lobby.ParticipantRepository) *JoinLobbyCommandHandler {
	return &JoinLobbyCommandHandler{EventPublisher: publishes, InstanceRepository: instances, ParticipantRepository: participants, ParticipantFactory: lobby.ParticipantFactory{}}
}

func (h *JoinLobbyCommandHandler) Handle(ctx *Context, cmd JoinLobbyCommand) error {

	if ctx.Session.CanJoinLobby() == false {
		return errors.New("player has already joined lobby")
	}

	instance, err := h.InstanceRepository.QueryByID(ctx.Context, cmd.LobbyID)
	if err != nil {
		return err
	}

	participant, err := h.ParticipantFactory.Create(ctx.ClientID, ctx.Session.PlayerID, instance, lobby.ParticipantFactoryOptions{
		RoleID:     uuid.Nil,
		SlotIndex:  cmd.SlotIndex,
		DeckIndex:  cmd.DeckIndex,
		UseStamina: cmd.UseStamina,
		FromInvite: cmd.FromInvite,
	})
	if err != nil {
		return err
	}

	if err := h.ParticipantRepository.Create(ctx.Context, participant); err != nil {
		return err
	}

	h.EventPublisher.Notify(lobby.NewParticipantCreatedEvent(ctx.Context, participant.ClientID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	return nil

}
