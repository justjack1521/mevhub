package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/domain/game"
	"mevhub/internal/domain/lobby"
)

type GameInstanceWriter struct {
	EventPublisher          *mevent.Publisher
	LobbyInstanceRepository lobby.InstanceReadRepository
	GameInstanceFactory     *game.InstanceFactory
	GameInstanceRepository  game.InstanceWriteRepository
}

func NewGameInstanceWriter(publisher *mevent.Publisher, lobbies lobby.InstanceRepository, factory *game.InstanceFactory, repository game.InstanceWriteRepository) *GameInstanceWriter {
	var subscriber = &GameInstanceWriter{EventPublisher: publisher, LobbyInstanceRepository: lobbies, GameInstanceFactory: factory, GameInstanceRepository: repository}
	publisher.Subscribe(subscriber, lobby.InstanceStartedEvent{})
	return subscriber
}

func (s *GameInstanceWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.InstanceStartedEvent:
		if err := s.HandleStart(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *GameInstanceWriter) HandleStart(event lobby.InstanceStartedEvent) error {

	source, err := s.LobbyInstanceRepository.QueryByID(event.Context(), event.LobbyID())
	if err != nil {
		return err
	}

	instance, err := s.GameInstanceFactory.Create(source)
	if err != nil {
		return err
	}

	if err := s.GameInstanceRepository.Create(event.Context(), instance); err != nil {
		return err
	}

	s.EventPublisher.Notify(game.NewInstanceCreatedEvent(event.Context(), instance.SysID))

	return nil

}
