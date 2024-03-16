package extern

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protogame"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
)

type LobbyPlayerSummaryExternalRepository struct {
	source *GamePlayerSummaryExternalRepository
}

func NewLobbyPlayerSummaryExternalRepository(game services.MeviusGameServiceClient) *LobbyPlayerSummaryExternalRepository {
	return &LobbyPlayerSummaryExternalRepository{source: NewGamePlayerSummaryExternalRepository(game)}
}

func (r *LobbyPlayerSummaryExternalRepository) Query(ctx context.Context, id uuid.UUID, deck int) (lobby.PlayerSummary, error) {
	summary, err := r.source.Query(ctx, id, deck)
	if err != nil {
		return lobby.PlayerSummary{}, err
	}
	var result = lobby.PlayerSummary{
		Identity: lobby.PlayerIdentity{
			PlayerID:      summary.PlayerID,
			PlayerName:    summary.PlayerName,
			PlayerComment: summary.PlayerComment,
			PlayerLevel:   summary.PlayerLevel,
		},
		Loadout: lobby.PlayerLoadout{
			DeckIndex: summary.DeckIndex,
			JobCard: lobby.PlayerJobCardSummary{
				JobCardID:      summary.JobCard.JobCardID,
				SubJobIndex:    summary.JobCard.SubJobIndex,
				CrownLevel:     summary.JobCard.CrownLevel,
				OverBoostLevel: summary.JobCard.OverBoostLevel,
			},
			Weapon: lobby.PlayerWeaponSummary{
				WeaponID:        summary.Weapon.WeaponID,
				SubWeaponUnlock: summary.Weapon.SubWeaponUnlock,
			},
			AbilityCards: make([]lobby.PlayerAbilityCardSummary, len(summary.AbilityCards)),
		},
	}
	for index, value := range summary.AbilityCards {
		result.Loadout.AbilityCards[index] = lobby.PlayerAbilityCardSummary{
			AbilityCardID:    value.AbilityCardID,
			SlotIndex:        value.SlotIndex,
			AbilityCardLevel: value.AbilityCardLevel,
			AbilityLevel:     value.AbilityLevel,
			OverBoostLevel:   value.OverBoostLevel,
		}
	}
	return result, nil
}

type GamePlayerSummaryExternalRepository struct {
	game       services.MeviusGameServiceClient
	translator translate.GamePlayerSummaryTranslator
}

func NewGamePlayerSummaryExternalRepository(game services.MeviusGameServiceClient) *GamePlayerSummaryExternalRepository {
	return &GamePlayerSummaryExternalRepository{game: game, translator: translate.NewGamePlayerSummaryTranslator()}
}

func (r *GamePlayerSummaryExternalRepository) Query(ctx context.Context, id uuid.UUID, index int) (game.PlayerSummary, error) {
	loadout, err := r.game.GetMultiPlayerLoadout(ctx, &protogame.GetMultiPlayerLoadoutRequest{PlayerId: id.String(), DeckIndex: int32(index)})
	if err != nil {
		return game.PlayerSummary{}, err
	}
	return r.translator.Unmarshall(loadout.Primary)
}
