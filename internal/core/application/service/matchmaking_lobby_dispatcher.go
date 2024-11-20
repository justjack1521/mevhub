package service

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/factory"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
)

type LobbyMatchmakingDispatcher struct {
	EventPublisher          *mevent.Publisher
	QuestRepository         port.QuestRepository
	LobbyInstanceRepository port.LobbyInstanceReadRepository
	GameInstanceRepository  port.GameInstanceRepository
	GameInstanceFactory     *factory.GameInstanceFactory
}

func NewLobbyMatchmakingDispatcher(publisher *mevent.Publisher, quests port.QuestRepository, lobbies port.LobbyInstanceReadRepository, games port.GameInstanceRepository, factory *factory.GameInstanceFactory) *LobbyMatchmakingDispatcher {
	return &LobbyMatchmakingDispatcher{EventPublisher: publisher, QuestRepository: quests, LobbyInstanceRepository: lobbies, GameInstanceRepository: games, GameInstanceFactory: factory}
}

func (d *LobbyMatchmakingDispatcher) Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, lobbies match.LobbyQueueEntryCollection) error {

	quest, err := d.QuestRepository.QueryByID(id)
	if err != nil {
		return err
	}

	if len(lobbies) < quest.Tier.GameMode.MaxLobbies {
		return match.ErrNotEnoughLobbiesForGame(len(lobbies), quest.Tier.GameMode.MaxLobbies)
	}

	var gameID = uuid.NewV4()

	var instances = make([]*lobby.Instance, len(lobbies))

	for index, value := range lobbies {
		instance, err := d.LobbyInstanceRepository.QueryByID(ctx, value.LobbyID)
		if err != nil {
			return err
		}
		instances[index] = instance
	}

	result, err := d.GameInstanceFactory.Create(gameID, instances...)
	if err != nil {
		return err
	}

	if err := d.GameInstanceRepository.Create(ctx, result); err != nil {
		return err
	}

	for _, value := range lobbies {
		d.EventPublisher.Notify(lobby.NewInstanceStartedEvent(ctx, value.LobbyID, result.SysID))
	}

	d.EventPublisher.Notify(game.NewInstanceCreatedEvent(ctx, result.SysID))

	return nil

}
