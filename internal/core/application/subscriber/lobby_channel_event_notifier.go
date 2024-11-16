package subscriber

import (
	"context"
	"fmt"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/core/application/consumer"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/port"
)

type LobbyChannelEventNotifier struct {
	EventPublisher          *mevent.Publisher
	PlayerSummaryRepository port.LobbyPlayerSummaryReadRepository
	PlayerSummaryTranslator translate.LobbyPlayerSummaryTranslator
}

type ClientNotification interface {
	MarshallBinary() ([]byte, error)
}

func NewLobbyChannelEventNotifier(publisher *mevent.Publisher, summary port.LobbyPlayerSummaryReadRepository, translator translate.LobbyPlayerSummaryTranslator) *LobbyChannelEventNotifier {
	var subscriber = &LobbyChannelEventNotifier{EventPublisher: publisher, PlayerSummaryRepository: summary, PlayerSummaryTranslator: translator}
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{}, lobby.ParticipantDeletedEvent{}, lobby.ParticipantReadyEvent{}, lobby.ParticipantUnreadyEvent{}, lobby.InstanceStartedEvent{}, game.InstanceReadyEvent{})
	return subscriber
}

func (s *LobbyChannelEventNotifier) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case lobby.ParticipantCreatedEvent:
		if err := s.HandleParticipantAdded(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantDeletedEvent:
		if err := s.HandleParticipantDeleted(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantReadyEvent:
		if err := s.HandleParticipantReady(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantUnreadyEvent:
		if err := s.HandleParticipantUnready(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.ParticipantDeckChangeEvent:
		if err := s.HandleParticipantDeckChange(actual); err != nil {
			fmt.Println(err)
		}
	case lobby.InstanceStartedEvent:
		if err := s.HandleLobbyStartEvent(actual); err != nil {
			fmt.Println(err)
		}
	case game.InstanceReadyEvent:
		if err := s.HandleGameReadyEvent(actual); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *LobbyChannelEventNotifier) HandleGameReadyEvent(event game.InstanceReadyEvent) error {

	var notification = &protomulti.LobbyReadyNotification{
		LobbyId: event.InstanceID().String(),
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_LOBBY_READY, event.InstanceID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleLobbyStartEvent(event lobby.InstanceStartedEvent) error {

	var notification = &protomulti.LobbyStartNotification{
		LobbyId: event.LobbyID().String(),
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_LOBBY_START, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantDeleted(event lobby.ParticipantDeletedEvent) error {

	var notification = &protomulti.ParticipantLeaveNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerId:   event.PlayerID().String(),
		PlayerSlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_PARTICIPANT_LEAVE, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantAdded(event lobby.ParticipantCreatedEvent) error {

	if uuid.Equal(event.PlayerID(), uuid.Nil) {
		return nil
	}

	summary, err := s.PlayerSummaryRepository.Query(event.Context(), event.PlayerID())
	if err != nil {
		return err
	}

	player, err := s.PlayerSummaryTranslator.Marshall(summary)
	if err != nil {
		return err
	}

	var notification = &protomulti.ParticipantJoinNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerId:   event.PlayerID().String(),
		DeckIndex:  int32(event.DeckIndex()),
		PlayerSlot: int32(event.SlotIndex()),
		Player:     player,
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_PARTICIPANT_JOIN, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantDeckChange(event lobby.ParticipantDeckChangeEvent) error {

	summary, err := s.PlayerSummaryRepository.Query(event.Context(), event.PlayerID())
	if err != nil {
		return err
	}

	player, err := s.PlayerSummaryTranslator.Marshall(summary)
	if err != nil {
		return err
	}

	var notification = &protomulti.ParticipantDeckChangeNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerSlot: int32(event.SlotIndex()),
		Player:     player,
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_PARTICIPANT_DECK_CHANGE, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantReady(event lobby.ParticipantReadyEvent) error {

	var notification = &protomulti.ParticipantReadyNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerSlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_PARTICIPANT_READY, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantUnready(event lobby.ParticipantUnreadyEvent) error {

	var notification = &protomulti.ParticipantUnreadyNotification{
		LobbyId:   event.LobbyID().String(),
		PartySlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiLobbyNotificationType_PARTICIPANT_UNREADY, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) publish(ctx context.Context, t protomulti.MultiLobbyNotificationType, id uuid.UUID, n ClientNotification) error {
	bytes, err := n.MarshallBinary()
	if err != nil {
		return err
	}
	var message = consumer.NewLobbyClientNotificationEvent(ctx, t, id, bytes)
	s.EventPublisher.Notify(message)
	return nil
}
