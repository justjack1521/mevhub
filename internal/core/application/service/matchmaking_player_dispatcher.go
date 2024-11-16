package service

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type PlayerMatchmakingDispatcher struct {
	EventPublisher            *mevent.Publisher
	SessionInstanceRepository session.InstanceReadRepository
	LobbyInstanceRepository   port.LobbyInstanceReadRepository
	ParticipantRepository     lobby.ParticipantRepository
}

func (s PlayerMatchmakingDispatcher) Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, entry match.LobbyQueueEntry, player match.PlayerQueueEntry) error {

	sesh, err := s.SessionInstanceRepository.QueryByID(ctx, player.UserID)
	if err != nil {
		return err
	}

	instance, err := s.LobbyInstanceRepository.QueryByID(ctx, entry.LobbyID)
	if err != nil {
		return err
	}

	count, err := s.ParticipantRepository.QueryCountForLobby(ctx, entry.LobbyID)
	if err != nil {
		return err
	}

	participant, err := s.ParticipantRepository.QueryParticipantForLobby(ctx, instance.SysID, count-1)
	if err != nil {
		return err
	}

	var options = lobby.ParticipantJoinOptions{
		RoleID:     uuid.UUID{},
		SlotIndex:  participant.PlayerSlot,
		DeckIndex:  sesh.DeckIndex,
		UseStamina: false,
		FromInvite: false,
	}

	if err := participant.SetPlayer(sesh.UserID, sesh.PlayerID, options); err != nil {
		return err
	}

	if err := s.ParticipantRepository.Create(ctx, participant); err != nil {
		return err
	}

	s.EventPublisher.Notify(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	return nil

}
