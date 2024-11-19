package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type LobbyInstanceWriter struct {
	EventPublisher          *mevent.Publisher
	LobbyInstanceRepository port.LobbyInstanceRepository
}

func NewLobbyInstanceWriter(publisher *mevent.Publisher, lobbies port.LobbyInstanceRepository) *LobbyInstanceWriter {
	var service = &LobbyInstanceWriter{EventPublisher: publisher, LobbyInstanceRepository: lobbies}
	publisher.Subscribe(service, session.InstanceDeletedEvent{})
	return service
}

func (s *LobbyInstanceWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case session.InstanceDeletedEvent:
		if err := s.HandleSessionDeleted(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyInstanceWriter) HandleSessionDeleted(evt session.InstanceDeletedEvent) error {

	instance, err := s.LobbyInstanceRepository.QueryByID(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	if instance.HostPlayerID != evt.PlayerID() {
		return nil
	}

	if err := s.LobbyInstanceRepository.Delete(evt.Context(), instance.SysID); err != nil {
		return err
	}

	return nil

}
