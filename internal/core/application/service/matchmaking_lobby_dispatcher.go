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
	ParticipantRepository   port.LobbyParticipantReadRepository
	SessionRepository       port.SessionInstanceRepository
	GameInstanceRepository  port.GameInstanceRepository
	GameInstanceFactory     *factory.GameInstanceFactory
}

func NewLobbyMatchmakingDispatcher(publisher *mevent.Publisher, quests port.QuestRepository, lobbies port.LobbyInstanceReadRepository, participants port.LobbyParticipantReadRepository, sessions port.SessionInstanceRepository, games port.GameInstanceRepository, factory *factory.GameInstanceFactory) *LobbyMatchmakingDispatcher {
	return &LobbyMatchmakingDispatcher{EventPublisher: publisher, QuestRepository: quests, LobbyInstanceRepository: lobbies, ParticipantRepository: participants, SessionRepository: sessions, GameInstanceRepository: games, GameInstanceFactory: factory}
}

func (d *LobbyMatchmakingDispatcher) Dispatch(ctx context.Context, mode game.ModeIdentifier, q uuid.UUID, lobbies match.LobbyQueueEntryCollection) error {

	quest, err := d.QuestRepository.QueryByID(q)
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

	// Move every participant's session from the lobby to the new game,
	// mirroring lobby_start: without this, matchmade players get no reconnect
	// grace period and are never removed from the live game on disconnect.
	for _, instance := range instances {
		participants, err := d.ParticipantRepository.QueryAllForLobby(ctx, instance.SysID)
		if err != nil {
			return err
		}
		for _, participant := range participants {
			if !participant.HasPlayer() {
				continue
			}
			session, err := d.SessionRepository.QueryByID(ctx, participant.UserID)
			if err != nil {
				continue
			}
			if session.LobbyID != instance.SysID {
				continue
			}
			session.GameID = result.SysID
			session.LobbyID = uuid.Nil
			session.PartySlot = 0
			if err := d.SessionRepository.Update(ctx, session); err != nil {
				return err
			}
		}
	}

	for _, value := range lobbies {
		d.EventPublisher.Notify(lobby.NewInstanceStartedEvent(ctx, value.LobbyID, result.SysID))
	}

	d.EventPublisher.Notify(game.NewInstanceCreatedEvent(ctx, result.SysID))

	return nil

}
