package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/match"
	"mevhub/internal/core/port"
	"time"
)

type LobbyPlayerQueueWriter struct {
	QueueRepository         port.MatchLobbyPlayerQueueWriteRepository
	QuestRepository         port.QuestRepository
	ParticipantRepository   port.LobbyParticipantReadRepository
	PlayerSummaryRepository port.LobbyPlayerSummaryReadRepository
}

func NewLobbyPlayerQueueWriter(publisher *mevent.Publisher, queues port.MatchLobbyPlayerQueueWriteRepository, quests port.QuestRepository, participants port.LobbyParticipantReadRepository, players port.LobbyPlayerSummaryReadRepository) *LobbyPlayerQueueWriter {
	var subscriber = &LobbyPlayerQueueWriter{QueueRepository: queues, QuestRepository: quests, ParticipantRepository: participants, PlayerSummaryRepository: players}
	publisher.Subscribe(subscriber, lobby.InstanceCreatedEvent{})
	return subscriber
}

func (s *LobbyPlayerQueueWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceCreatedEvent:
		if err := s.Handle(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyPlayerQueueWriter) Handle(evt lobby.InstanceCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return nil
	}

	participants, err := s.ParticipantRepository.QueryAllForLobby(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	var filled int
	for _, participant := range participants {
		if participant.HasPlayer() {
			filled++
		}
	}

	if filled >= quest.Tier.GameMode.MaxPlayers {
		return nil
	}

	var sum int
	for _, participant := range participants {
		if participant.HasPlayer() == false {
			continue
		}
		player, err := s.PlayerSummaryRepository.Query(evt.Context(), participant.PlayerID)
		if err != nil {
			return err
		}
		sum += player.Loadout.CalculateDeckLevel()
	}
	var average = sum / filled

	var entry = match.LobbyQueueEntry{
		LobbyID:      evt.LobbyID(),
		QuestID:      evt.QuestID(),
		AverageLevel: average,
		JoinedAt:     time.Now().UTC(),
	}

	if err := s.QueueRepository.AddLobbyToQueue(evt.Context(), quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
		return err
	}

	return nil

}
