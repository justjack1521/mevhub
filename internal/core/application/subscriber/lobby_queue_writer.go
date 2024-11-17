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

type LobbyQueueWriter struct {
	LobbyRepository       port.LobbyInstanceReadRepository
	QueueRepository       port.MatchPlayerQueueWriteRepository
	QuestRepository       port.QuestRepository
	ParticipantRepository port.LobbyParticipantReadRepository
}

func NewLobbyQueueWriter(publisher *mevent.Publisher, lobbies port.LobbyInstanceReadRepository, queues port.MatchPlayerQueueWriteRepository, quests port.QuestRepository, participants port.LobbyParticipantReadRepository) *LobbyQueueWriter {
	var subscriber = &LobbyQueueWriter{LobbyRepository: lobbies, QueueRepository: queues, QuestRepository: quests, ParticipantRepository: participants}
	publisher.Subscribe(subscriber, lobby.InstanceCreatedEvent{})
	return subscriber
}

func (s *LobbyQueueWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceCreatedEvent:
		if err := s.Handle(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyQueueWriter) Handle(evt lobby.InstanceCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodMatch {
		return nil
	}

	instance, err := s.LobbyRepository.QueryByID(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	participants, err := s.ParticipantRepository.QueryAllForLobby(evt.Context(), instance.SysID)
	if err != nil {
		return err
	}

	if len(participants) == quest.Tier.GameMode.MaxPlayers {
		return nil
	}

	var sum int
	for _, participant := range participants {
		sum += participant.PlayerSlot * 10
	}
	var average = sum / len(participants)

	var entry = match.LobbyQueueEntry{
		LobbyID:      instance.SysID,
		QuestID:      instance.QuestID,
		AverageLevel: average,
		JoinedAt:     time.Now().UTC(),
	}

	if err := s.QueueRepository.AddLobbyToQueue(evt.Context(), quest.Tier.GameMode.ModeIdentifier, entry); err != nil {
		return err
	}

	return nil

}
