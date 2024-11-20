package external

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/game"
)

type GamePlayerLoadoutRepository struct {
	client     services.MeviusIdentityServiceClient
	translator translate.GamePlayerLoadoutTranslator
}

func NewGamePlayerLoadoutRepository(client services.MeviusIdentityServiceClient) *GamePlayerLoadoutRepository {
	return &GamePlayerLoadoutRepository{client: client, translator: translate.NewGamePlayerLoadoutTranslator()}
}

func (r *GamePlayerLoadoutRepository) Query(ctx context.Context, id uuid.UUID, index int) (game.PlayerLoadout, error) {
	loadout, err := r.client.GetSinglePlayerLoadout(ctx, &protoidentity.GetSinglePlayerLoadoutRequest{
		PlayerId:  id.String(),
		DeckIndex: int32(index),
	})
	if err != nil {
		return game.PlayerLoadout{}, err
	}
	translated, err := r.translator.Unmarshall(loadout.Loadout)
	if err != nil {
		return game.PlayerLoadout{}, err
	}
	return translated, nil
}
