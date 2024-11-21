package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

var (
	lobbySummaryCreateFailed = func(err error) error {
		return fmt.Errorf("failed to create lobby summary: %w", err)
	}
	lobbySummaryDeleteFailed = func(err error) error {
		return fmt.Errorf("failed to delete lobby summary: %w", err)
	}
)

type LobbySummaryWriter struct {
	EventPublisher    *mevent.Publisher
	QuestRepository   port.QuestRepository
	SummaryRepository port.LobbySummaryWriteRepository
}

func NewLobbySummaryWriter(publisher *mevent.Publisher, quests port.QuestRepository, summaries port.LobbySummaryWriteRepository) *LobbySummaryWriter {
	var subscriber = &LobbySummaryWriter{EventPublisher: publisher, QuestRepository: quests, SummaryRepository: summaries}
	publisher.Subscribe(subscriber, lobby.InstanceCreatedEvent{}, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbySummaryWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceCreatedEvent:
		if err := s.HandleCreate(actual); err != nil {
			fmt.Println(lobbySummaryCreateFailed(err))
		}
	case lobby.InstanceDeletedEvent:
		if err := s.HandleDelete(actual); err != nil {
			fmt.Println(lobbySummaryDeleteFailed(err))
		}
	}
}

func (s *LobbySummaryWriter) HandleCreate(evt lobby.InstanceCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return port.ErrFailedGetQuestByID(evt.QuestID(), err)
	}

	var summary = lobby.Summary{
		InstanceID:         evt.LobbyID(),
		QuestID:            quest.SysID,
		PartyID:            evt.PartyID(),
		LobbyComment:       evt.Comment(),
		MinimumPlayerLevel: evt.MinPlayerLevel(),
	}

	if err := s.SummaryRepository.Create(evt.Context(), summary); err != nil {
		return port.ErrFailedCreateLobbySummary(summary, err)
	}

	var categories = make([]uuid.UUID, len(quest.Categories))
	for index, value := range quest.Categories {
		categories[index] = value.SysID
	}

	s.EventPublisher.Notify(lobby.NewSummaryCreatedEvent(evt.Context(), summary.InstanceID, quest.SysID, quest.Tier.StarLevel, summary.MinimumPlayerLevel, categories))

	return nil

}

func (s *LobbySummaryWriter) HandleDelete(evt lobby.InstanceDeletedEvent) error {
	if err := s.SummaryRepository.Delete(evt.Context(), evt.LobbyID()); err != nil {
		return port.ErrFailedDeleteLobbySummary(evt.LobbyID(), err)
	}
	return nil
}
