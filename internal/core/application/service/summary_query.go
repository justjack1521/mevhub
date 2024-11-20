package service

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

var (
	ErrFailedQueryLobbyPartySummary = func(party string, err error) error {
		return fmt.Errorf("failed to query lobby party %s summary: %w", party, err)
	}
	ErrFailedQueryLobbySummary = func(id uuid.UUID, err error) error {
		return fmt.Errorf("failed to query lobby %s summary: %w", id, err)
	}
)

type SummaryQueryService struct {
	InstanceRepository    port.LobbyInstanceReadRepository
	ParticipantRepository port.LobbyParticipantReadRepository
	LobbySummary          port.LobbySummaryReadRepository
	PlayerSummary         port.LobbyPlayerSummaryReadRepository
}

func NewSummaryQueryService(instances port.LobbyInstanceReadRepository, participants port.LobbyParticipantReadRepository, lobbies port.LobbySummaryReadRepository, players port.LobbyPlayerSummaryReadRepository) *SummaryQueryService {
	return &SummaryQueryService{InstanceRepository: instances, ParticipantRepository: participants, LobbySummary: lobbies, PlayerSummary: players}
}

func (s *SummaryQueryService) QueryByPartyID(ctx context.Context, party string) (lobby.Summary, error) {
	instance, err := s.InstanceRepository.QueryByPartyID(ctx, party)
	if err != nil {
		return lobby.Summary{}, ErrFailedQueryLobbyPartySummary(party, err)
	}
	summary, err := s.Query(ctx, instance.SysID)
	if err != nil {
		return lobby.Summary{}, ErrFailedQueryLobbyPartySummary(party, err)
	}
	return summary, nil
}

func (s *SummaryQueryService) Query(ctx context.Context, id uuid.UUID) (lobby.Summary, error) {

	instance, err := s.InstanceRepository.QueryByID(ctx, id)
	if err != nil {
		return lobby.Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	summary, err := s.LobbySummary.Query(ctx, instance.SysID)
	if err != nil {
		return lobby.Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	participants, err := s.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return lobby.Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	var players = make([]lobby.PlayerSlotSummary, len(participants))

	for index, value := range participants {
		if uuid.Equal(value.PlayerID, uuid.Nil) {
			players[index] = lobby.PlayerSlotSummary{
				PartySlot:     value.PlayerSlot,
				Ready:         false,
				PlayerSummary: lobby.PlayerSummary{},
			}
			continue
		}
		player, err := s.PlayerSummary.Query(ctx, value.PlayerID)
		if err != nil {
			return lobby.Summary{}, ErrFailedQueryLobbySummary(id, err)
		}
		players[index] = lobby.PlayerSlotSummary{
			PartySlot:     value.PlayerSlot,
			Ready:         value.Ready,
			PlayerSummary: player,
		}
	}

	summary.Players = players

	return summary, nil

}
