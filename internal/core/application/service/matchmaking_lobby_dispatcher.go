package service

import (
	"context"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
)

type LobbyMatchmakingDispatcher struct {
	EventPublisher          *mevent.Publisher
	QuestRepository         port.QuestRepository
	LobbyInstanceRepository port.LobbyInstanceReadRepository
	ParticipantRepository   port.LobbyParticipantRepository
}

func NewLobbyMatchmakingDispatcher(publisher *mevent.Publisher, quests port.QuestRepository, lobbies port.LobbyInstanceReadRepository, participants port.LobbyParticipantRepository) *LobbyMatchmakingDispatcher {
	return &LobbyMatchmakingDispatcher{EventPublisher: publisher, QuestRepository: quests, LobbyInstanceRepository: lobbies, ParticipantRepository: participants}
}

func (d *LobbyMatchmakingDispatcher) Dispatch(ctx context.Context, mode game.ModeIdentifier, id uuid.UUID, lobbies match.LobbyQueueEntryCollection) error {

	quest, err := d.QuestRepository.QueryByID(id)
	if err != nil {
		return err
	}

	if len(lobbies) < quest.Tier.GameMode.MaxLobbies {
		return match.ErrNotEnoughLobbiesForGame(len(lobbies), quest.Tier.GameMode.MaxLobbies)
	}

	return nil

}
