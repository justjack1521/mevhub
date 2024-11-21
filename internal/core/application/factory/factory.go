package factory

import (
	uuid "github.com/satori/go.uuid"
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

func (f *GameInstanceFactory) Create(id uuid.UUID, sources ...*lobby.Instance) (*game.Instance, error) {

	var ids = make([]uuid.UUID, len(sources))
	for index, value := range sources {
		ids[index] = value.SysID
	}

	quest, err := f.QuestRepository.QueryByID(sources[0].QuestID)
	if err != nil {
		return nil, err
	}

	return &game.Instance{
		SysID:    id,
		Seed:     rand.Int(),
		LobbyIDs: ids,
		Options: &game.InstanceOptions{
			MinimumPlayerLevel: sources[0].MinimumPlayerLevel,
			MaxRunTime:         quest.Tier.TimeLimit,
			PlayerTurnDuration: quest.Tier.PlayerTurnDuration,
			MaxPartyCount:      quest.Tier.GameMode.MaxLobbies,
			MaxPlayerCount:     quest.Tier.GameMode.MaxPlayers,
		},
		State:        game.InstanceGamePendingState,
		RegisteredAt: time.Now(),
	}, nil
}
