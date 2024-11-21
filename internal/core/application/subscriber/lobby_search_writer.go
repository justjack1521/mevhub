package subscriber

import (
	"fmt"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbySearchWriter struct {
	LobbyRepository  port.LobbyInstanceReadRepository
	SearchRepository port.LobbySearchWriteRepository
	QuestRepository  port.QuestRepository
}

func NewLobbySearchWriter(publisher *mevent.Publisher, lobbies port.LobbyInstanceReadRepository, search port.LobbySearchWriteRepository, quests port.QuestRepository) *LobbySearchWriter {
	var subscriber = &LobbySearchWriter{LobbyRepository: lobbies, SearchRepository: search, QuestRepository: quests}
	publisher.Subscribe(subscriber, lobby.InstanceCreatedEvent{})
	return subscriber
}

func (s *LobbySearchWriter) Notify(event mevent.Event) {
	actual, valid := event.(lobby.InstanceCreatedEvent)
	if valid == false {
		return
	}
	if err := s.Handle(actual); err != nil {
		fmt.Println(err)
	}
}

func (s *LobbySearchWriter) Handle(evt lobby.InstanceCreatedEvent) error {

	quest, err := s.QuestRepository.QueryByID(evt.QuestID())
	if err != nil {
		return err
	}

	if quest.Tier.GameMode.FulfillMethod != game.FulfillMethodSearch {
		return nil
	}

	instance, err := s.LobbyRepository.QueryByID(evt.Context(), evt.LobbyID())
	if err != nil {
		return err
	}

	var categories = make([]uuid.UUID, len(quest.Categories))
	for index, category := range quest.Categories {
		if category.Zero() {
			continue
		}
		categories[index] = category.SysID
	}

	var search = lobby.SearchEntry{
		InstanceID:         evt.LobbyID(),
		ModeIdentifier:     string(quest.Tier.GameMode.ModeIdentifier),
		Level:              quest.Tier.StarLevel,
		MinimumPlayerLevel: instance.MinimumPlayerLevel,
		Categories:         categories,
	}

	if err := s.SearchRepository.Create(evt.Context(), search); err != nil {
		return err
	}

	return nil

}
