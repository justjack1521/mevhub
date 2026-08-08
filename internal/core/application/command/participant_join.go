package command

import (
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type ParticipantJoinCommand struct {
	BasicCommand
	LobbyID    uuid.UUID
	DeckIndex  int
	SlotIndex  int
	UseStamina bool
	FromInvite bool
}

func (c ParticipantJoinCommand) CommandName() string {
	return "participant.lobby"
}

func NewParticipantJoinCommand(id uuid.UUID, deck, slot int, stamina, invite bool) *ParticipantJoinCommand {
	return &ParticipantJoinCommand{
		LobbyID:    id,
		DeckIndex:  deck,
		SlotIndex:  slot,
		UseStamina: stamina,
		FromInvite: invite,
	}
}

type ParticipantJoinCommandHandler struct {
	EventPublisher        *mevent.Publisher
	SessionRepository     port.SessionInstanceRepository
	InstanceRepository    port.LobbyInstanceReadRepository
	ParticipantFactory    lobby.ParticipantFactory
	ParticipantRepository port.LobbyParticipantRepository
	ListenerRepository    lobby.NotificationListenerWriteRepository
}

func NewParticipantJoinCommandHandler(publishes *mevent.Publisher, sessions port.SessionInstanceRepository, instances port.LobbyInstanceReadRepository, participants port.LobbyParticipantRepository, listeners lobby.NotificationListenerWriteRepository) *ParticipantJoinCommandHandler {
	return &ParticipantJoinCommandHandler{EventPublisher: publishes, SessionRepository: sessions, InstanceRepository: instances, ParticipantRepository: participants, ParticipantFactory: lobby.ParticipantFactory{}, ListenerRepository: listeners}
}

func (h *ParticipantJoinCommandHandler) Handle(ctx Context, cmd *ParticipantJoinCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}
	if !current.CanJoinLobby() {
		return lobby.ErrPlayerAlreadyInLobby(current.LobbyID)
	}

	instance, err := h.InstanceRepository.QueryByID(ctx, cmd.LobbyID)
	if err != nil {
		return err
	}

	participant, err := h.ParticipantRepository.QueryParticipantForLobby(ctx, instance.SysID, cmd.SlotIndex)
	if err != nil {
		return err
	}

	if participant.HasPlayer() {
		return lobby.ErrPlayerSlotAlreadyTaken(cmd.SlotIndex)
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

	current.LobbyID = participant.LobbyID
	current.PartySlot = participant.PlayerSlot
	if err := h.SessionRepository.Update(ctx, current); err != nil {
		return err
	}

	if err := h.ListenerRepository.CreateListener(ctx, participant.LobbyID, ctx.UserID(), ctx.PlayerID()); err != nil {
		return err
	}

	cmd.QueueEvent(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	return nil

}
