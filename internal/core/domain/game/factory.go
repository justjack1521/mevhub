package game

import (
	"math/rand"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
	"time"
)

type InstanceFactory struct {
	QuestRepository port.QuestRepository
}

func NewInstanceFactory(quests port.QuestRepository) *InstanceFactory {
	return &InstanceFactory{QuestRepository: quests}
}

func (f *InstanceFactory) Create(source *lobby.Instance) (*Instance, error) {

	quest, err := f.QuestRepository.QueryByID(source.QuestID)
	if err != nil {
		return nil, err
	}

	return &Instance{
		SysID: source.SysID,
		Seed:  rand.Int(),
		Options: &InstanceOptions{
			MinimumPlayerLevel: source.MinimumPlayerLevel,
			MaxRunTime:         quest.Tier.TimeLimit,
			PlayerTurnDuration: quest.Tier.PlayerTurnDuration,
			MaxPlayerCount:     quest.Tier.GameMode.MaxPlayers,
		},
		State:        InstanceGamePendingState,
		RegisteredAt: time.Now(),
	}, nil
}
