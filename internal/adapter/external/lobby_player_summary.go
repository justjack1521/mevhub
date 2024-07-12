package external

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/domain/lobby"
)

type LobbyPlayerSummaryRepository struct {
	client     services.MeviusIdentityServiceClient
	translator translate.LobbyPlayerSummaryTranslator
}

func NewLobbyPlayerSummaryRepository(client services.MeviusIdentityServiceClient) *LobbyPlayerSummaryRepository {
	return &LobbyPlayerSummaryRepository{client: client, translator: translate.NewLobbyPlayerSummaryTranslator()}
}

func (r *LobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID) (lobby.PlayerSummary, error) {
	identity, err := r.client.GetSinglePlayerLoadoutIdentity(ctx, &protoidentity.GetSinglePlayerLoadoutIdentityRequest{PlayerId: id.String()})
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	translated, err := r.translator.Unmarshall(&protomulti.ProtoLobbyPlayer{
		Identity: identity.Identity,
		Loadout:  identity.Loadout,
	})
	if err != nil {
		return lobby.PlayerSummary{}, err
	}

	return translated, nil

}
