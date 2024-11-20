package factory

import (
	"context"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/port"
)

type GamePlayerFactory struct {
	loadout port.GamePlayerLoadoutReadRepository
}

func NewGamePlayerFactory(loadout port.GamePlayerLoadoutReadRepository) *GamePlayerFactory {
	return &GamePlayerFactory{loadout: loadout}
}

func (f *GamePlayerFactory) Create(ctx context.Context, source *game.Participant) (game.Player, error) {

	loadout, err := f.loadout.Query(ctx, source.PlayerID, source.DeckIndex)
	if err != nil {
		return game.Player{}, err
	}

	return game.Player{
		UserID:     source.UserID,
		PlayerID:   source.PlayerID,
		PlayerSlot: source.PlayerSlot,
		BotControl: source.BotControl,
		Loadout:    loadout,
	}, nil

}
