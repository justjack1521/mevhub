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
	QuestRepository       port.QuestRepository
	QueueRepository       port.MatchLobbyQueueWriteRepository
	ParticipantRepository port.LobbyParticipantReadRepository
}

func NewLobbyQueueWriter(publisher *mevent.Publisher, instance port.LobbyInstanceReadRepository, quests port.QuestRepository, queue port.MatchLobbyQueueWriteRepository, participants port.LobbyParticipantReadRepository) *LobbyQueueWriter {
	var subscriber = &LobbyQueueWriter{LobbyRepository: instance, QuestRepository: quests, QueueRepository: queue, ParticipantRepository: participants}
	publisher.Subscribe(subscriber, lobby.InstanceReadyEvent{})
	return subscriber
}

func (s *LobbyQueueWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceReadyEvent:
		if err := s.Handle(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyQueueWriter) Handle(evt lobby.InstanceReadyEvent) error {

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
