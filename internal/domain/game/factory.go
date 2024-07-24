package game

import (
	"math/rand"
	"mevhub/internal/domain/lobby"
	"time"
)

type InstanceFactory struct {
	QuestRepository QuestRepository
}

func NewInstanceFactory(quests QuestRepository) *InstanceFactory {
	return &InstanceFactory{QuestRepository: quests}
}

func (f *InstanceFactory) Create(source *lobby.Instance) (*Instance, error) {

	quest, err := f.QuestRepository.QueryByID(source.QuestID)
	if err != nil {
		return nil, err
	}

	return &Instance{
		SysID:   source.SysID,
		PartyID: source.PartyID,
		Seed:    rand.Int(),
		Options: &InstanceOptions{
			MinimumPlayerLevel: source.MinimumPlayerLevel,
			MaxRunTime:         quest.Tier.TimeLimit,
			PlayerTurnDuration: time.Second * 30, //quest.Tier.PlayerTurnDuration,
		},
		State:        InstanceGamePendingState,
		RegisteredAt: time.Now(),
	}, nil
}
