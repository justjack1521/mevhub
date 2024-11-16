package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbySummaryWriter struct {
	EventPublisher    *mevent.Publisher
	QuestRepository   port.QuestRepository
	SummaryRepository lobby.SummaryWriteRepository
}

func NewLobbySummaryWriter(publisher *mevent.Publisher, quests port.QuestRepository, summaries lobby.SummaryWriteRepository) *LobbySummaryWriter {
	var subscriber = &LobbySummaryWriter{EventPublisher: publisher, QuestRepository: quests, SummaryRepository: summaries}
	publisher.Subscribe(subscriber, lobby.InstanceCreatedEvent{}, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbySummaryWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceCreatedEvent:
		if err := s.HandleCreate(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.InstanceDeletedEvent:
		if err := s.HandleDelete(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbySummaryWriter) HandleDelete(event lobby.InstanceDeletedEvent) error {
	return s.SummaryRepository.Delete(event.Context(), event.LobbyID())
}

func (s *LobbySummaryWriter) HandleCreate(event lobby.InstanceCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(event.QuestID())
	if err != nil {
		return err
	}

	var summary = lobby.Summary{
		InstanceID:         event.LobbyID(),
		QuestID:            quest.SysID,
		PartyID:            event.PartyID(),
		LobbyComment:       event.Comment(),
		MinimumPlayerLevel: event.MinPlayerLevel(),
	}

	if err := s.SummaryRepository.Create(event.Context(), summary); err != nil {
		return err
	}

	var categories = make([]uuid.UUID, len(quest.Categories))
	for index, value := range quest.Categories {
		categories[index] = value.SysID
	}

	var evt = lobby.NewSummaryCreatedEvent(event.Context(), summary.InstanceID, quest.SysID, quest.Tier.StarLevel, summary.MinimumPlayerLevel, categories)
	s.EventPublisher.Notify(evt)

	return nil

}
