package lobby

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
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
	ParticipantRepository ParticipantReadRepository
	LobbySummary          SummaryReadRepository
	PlayerSummary         PlayerSummaryReadRepository
}

func NewSummaryQueryService(instances port.LobbyInstanceReadRepository, participants ParticipantReadRepository, lobbies SummaryReadRepository, players PlayerSummaryReadRepository) *SummaryQueryService {
	return &SummaryQueryService{InstanceRepository: instances, ParticipantRepository: participants, LobbySummary: lobbies, PlayerSummary: players}
}

func (s *SummaryQueryService) QueryByPartyID(ctx context.Context, party string) (Summary, error) {
	instance, err := s.InstanceRepository.QueryByPartyID(ctx, party)
	if err != nil {
		return Summary{}, ErrFailedQueryLobbyPartySummary(party, err)
	}
	summary, err := s.QueryByID(ctx, instance.SysID)
	if err != nil {
		return Summary{}, ErrFailedQueryLobbyPartySummary(party, err)
	}
	return summary, nil
}

func (s *SummaryQueryService) QueryByID(ctx context.Context, id uuid.UUID) (Summary, error) {

	instance, err := s.InstanceRepository.QueryByID(ctx, id)
	if err != nil {
		return Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	summary, err := s.LobbySummary.QueryByID(ctx, instance.SysID)
	if err != nil {
		return Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	participants, err := s.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
	if err != nil {
		return Summary{}, ErrFailedQueryLobbySummary(id, err)
	}

	var players = make([]PlayerSlotSummary, len(participants))

	for index, value := range participants {
		if uuid.Equal(value.PlayerID, uuid.Nil) {
			players[index] = PlayerSlotSummary{
				PartySlot:     value.PlayerSlot,
				Ready:         false,
				PlayerSummary: PlayerSummary{},
			}
			continue
		}
		player, err := s.PlayerSummary.Query(ctx, value.PlayerID)
		if err != nil {
			return Summary{}, ErrFailedQueryLobbySummary(id, err)
		}
		players[index] = PlayerSlotSummary{
			PartySlot:     value.PlayerSlot,
			Ready:         value.Ready,
			PlayerSummary: player,
		}
	}

	summary.Players = players

	return summary, nil

}
