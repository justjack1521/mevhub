package external

import (
	"context"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/domain/lobby"

	"github.com/justjack1521/mevium/pkg/genproto/protoidentity"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	services "github.com/justjack1521/mevium/pkg/genproto/service"
	uuid "github.com/satori/go.uuid"
)

type LobbyPlayerSummaryRepository struct {
	client     services.MeviusIdentityServiceClient
	translator translate.LobbyPlayerSummaryTranslator
}

func NewLobbyPlayerSummaryRepository(client services.MeviusIdentityServiceClient) *LobbyPlayerSummaryRepository {
	return &LobbyPlayerSummaryRepository{client: client, translator: translate.NewLobbyPlayerSummaryTranslator()}
}

func (r *LobbyPlayerSummaryRepository) Query(ctx context.Context, id uuid.UUID, deck int) (lobby.PlayerSummary, error) {
	// The deck index was never sent, so mobius answered for deck 0 whatever the player had joined
	// with - the lobby showed the wrong deck while the game, which asks with the index, was right.
	identity, err := r.client.GetMultiPlayerLoadoutIdentity(ctx, &protoidentity.GetMultiPlayerLoadoutIdentityRequest{PlayerId: id.String(), DeckIndex: int32(deck)})
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
