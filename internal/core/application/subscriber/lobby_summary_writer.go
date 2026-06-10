package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

var (
	lobbySummaryDeleteFailed = func(err error) error {
		return fmt.Errorf("failed to delete lobby summary: %w", err)
	}
)

type LobbySummaryWriter struct {
	SummaryRepository port.LobbySummaryWriteRepository
}

func NewLobbySummaryWriter(publisher *mevent.Publisher, summaries port.LobbySummaryWriteRepository) *LobbySummaryWriter {
	var subscriber = &LobbySummaryWriter{SummaryRepository: summaries}
	publisher.Subscribe(subscriber, lobby.InstanceDeletedEvent{})
	return subscriber
}

func (s *LobbySummaryWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceDeletedEvent:
		if err := s.HandleDelete(actual); err != nil {
			fmt.Println(lobbySummaryDeleteFailed(err))
		}
	}
}

func (s *LobbySummaryWriter) HandleDelete(evt lobby.InstanceDeletedEvent) error {
	if err := s.SummaryRepository.Delete(evt.Context(), evt.LobbyID()); err != nil {
		return port.ErrFailedDeleteLobbySummary(evt.LobbyID(), err)
	}
	return nil
}
