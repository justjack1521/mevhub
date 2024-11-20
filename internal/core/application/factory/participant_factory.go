package factory

import (
	"context"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type PlayerParticipantFactory struct {
	loadout port.GamePlayerLoadoutReadRepository
}

func NewPlayerParticipantFactory(loadout port.GamePlayerLoadoutReadRepository) *PlayerParticipantFactory {
	return &PlayerParticipantFactory{loadout: loadout}
}

func (f *PlayerParticipantFactory) Create(ctx context.Context, source *lobby.Participant) (game.Player, error) {

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
