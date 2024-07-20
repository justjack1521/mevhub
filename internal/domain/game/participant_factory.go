package game

import (
	"context"
	"mevhub/internal/domain/lobby"
)

type PlayerParticipantFactory struct {
	loadout PlayerLoadoutReadRepository
}

func NewPlayerParticipantFactory(loadout PlayerLoadoutReadRepository) *PlayerParticipantFactory {
	return &PlayerParticipantFactory{loadout: loadout}
}

func (f *PlayerParticipantFactory) Create(ctx context.Context, source *lobby.Participant) (*PlayerParticipant, error) {

	loadout, err := f.loadout.Query(ctx, source.PlayerID, source.DeckIndex)
	if err != nil {
		return nil, err
	}

	return &PlayerParticipant{
		UserID:     source.UserID,
		PlayerID:   source.PlayerID,
		PlayerSlot: source.PlayerSlot,
		BotControl: source.BotControl,
		Loadout:    loadout,
	}, nil

}
