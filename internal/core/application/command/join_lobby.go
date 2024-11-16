package command

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type JoinLobbyCommand struct {
	BasicCommand
	LobbyID    uuid.UUID
	DeckIndex  int
	SlotIndex  int
	UseStamina bool
	FromInvite bool
}

func (c JoinLobbyCommand) CommandName() string {
	return "join.lobby"
}

func NewJoinLobbyCommand(id uuid.UUID, deck, slot int, stamina, invite bool) *JoinLobbyCommand {
	return &JoinLobbyCommand{
		LobbyID:    id,
		DeckIndex:  deck,
		SlotIndex:  slot,
		UseStamina: stamina,
		FromInvite: invite,
	}
}

type JoinLobbyCommandHandler struct {
	EventPublisher        *mevent.Publisher
	InstanceRepository    port.LobbyInstanceReadRepository
	ParticipantFactory    lobby.ParticipantFactory
	ParticipantRepository port.LobbyParticipantRepository
}

func NewJoinLobbyCommandHandler(publishes *mevent.Publisher, instances port.LobbyInstanceReadRepository, participants port.LobbyParticipantRepository) *JoinLobbyCommandHandler {
	return &JoinLobbyCommandHandler{EventPublisher: publishes, InstanceRepository: instances, ParticipantRepository: participants, ParticipantFactory: lobby.ParticipantFactory{}}
}

func (h *JoinLobbyCommandHandler) Handle(ctx Context, cmd *JoinLobbyCommand) error {

	instance, err := h.InstanceRepository.QueryByID(ctx, cmd.LobbyID)
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, instance.SysID, cmd.SlotIndex)
	if err != nil {
		return err
	}

	var options = lobby.ParticipantJoinOptions{
		RoleID:     uuid.UUID{},
		SlotIndex:  cmd.SlotIndex,
		DeckIndex:  cmd.DeckIndex,
		UseStamina: cmd.UseStamina,
		FromInvite: cmd.FromInvite,
	}

	if err := participant.SetPlayer(ctx.UserID(), ctx.PlayerID(), options); err != nil {
		return err
	}

	if err := h.ParticipantRepository.Create(ctx, participant); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	return nil

}
