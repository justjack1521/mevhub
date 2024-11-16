package factory

import (
	"math/rand"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
	"time"
)

type GameInstanceFactory struct {
	QuestRepository port.QuestRepository
}

func NewGameInstanceFactory(quests port.QuestRepository) *GameInstanceFactory {
	return &GameInstanceFactory{QuestRepository: quests}
}

func (f *GameInstanceFactory) Create(source *lobby.Instance) (*game.Instance, error) {

	quest, err := f.QuestRepository.QueryByID(source.QuestID)
	if err != nil {
		return nil, err
	}

	return &game.Instance{
		SysID: source.SysID,
		Seed:  rand.Int(),
		Options: &game.InstanceOptions{
			MinimumPlayerLevel: source.MinimumPlayerLevel,
			MaxRunTime:         quest.Tier.TimeLimit,
			PlayerTurnDuration: quest.Tier.PlayerTurnDuration,
			MaxPlayerCount:     quest.Tier.GameMode.MaxPlayers,
		},
		State:        game.InstanceGamePendingState,
		RegisteredAt: time.Now(),
	}, nil
}
