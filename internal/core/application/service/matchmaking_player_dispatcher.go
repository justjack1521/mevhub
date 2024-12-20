package service

import (
	"context"
	"errors"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
)

type PlayerMatchmakingDispatcher struct {
	EventPublisher            *mevent.Publisher
	SessionInstanceRepository port.SessionInstanceReadRepository
	LobbyInstanceRepository   port.LobbyInstanceReadRepository
	QuestRepository           port.QuestRepository
	ParticipantRepository     port.LobbyParticipantRepository
}

func NewPlayerMatchmakingDispatcher(publisher *mevent.Publisher, sessions port.SessionInstanceReadRepository, lobbies port.LobbyInstanceReadRepository, quests port.QuestRepository, participants port.LobbyParticipantRepository) *PlayerMatchmakingDispatcher {
	return &PlayerMatchmakingDispatcher{EventPublisher: publisher, SessionInstanceRepository: sessions, LobbyInstanceRepository: lobbies, ParticipantRepository: participants}
}

func (s PlayerMatchmakingDispatcher) Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, entry match.LobbyQueueEntry, player match.PlayerQueueEntry) (bool, error) {

	sesh, err := s.SessionInstanceRepository.QueryByID(ctx, player.UserID)
	if err != nil {
		return false, err
	}

	instance, err := s.LobbyInstanceRepository.QueryByID(ctx, entry.LobbyID)
	if err != nil {
		return false, err
	}

	existing, err := s.ParticipantRepository.QueryAllForLobby(ctx, entry.LobbyID)
	if err != nil {
		return false, err
	}

	quest, err := s.QuestRepository.QueryByID(entry.QuestID)
	if err != nil {
		return false, err
	}

	var filled int
	for _, exist := range existing {
		if exist.HasPlayer() {
			filled++
		}
	}

	if filled == quest.Tier.GameMode.MaxPlayers {
		return false, errors.New("lobby is full")
	}

	participant, err := s.ParticipantRepository.QueryParticipantForLobby(ctx, instance.SysID, filled)
	if err != nil {
		return false, err
	}

	var options = lobby.ParticipantJoinOptions{
		RoleID:     uuid.UUID{},
		SlotIndex:  participant.PlayerSlot,
		DeckIndex:  sesh.DeckIndex,
		UseStamina: false,
		FromInvite: false,
	}

	if err := participant.SetPlayer(sesh.UserID, sesh.PlayerID, options); err != nil {
		return false, err
	}

	if err := s.ParticipantRepository.Create(ctx, participant); err != nil {
		return false, err
	}

	s.EventPublisher.Notify(lobby.NewParticipantCreatedEvent(ctx, participant.UserID, participant.PlayerID, participant.LobbyID, participant.DeckIndex, participant.PlayerSlot))

	filled++

	return filled == instance.PlayerSlotCount, nil

}
