package subscriber

import (
	"context"
	"fmt"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/adapter/translate"
	"mevhub/internal/app/consumer"
	"mevhub/internal/domain/lobby"
)

type LobbyChannelEventNotifier struct {
	EventPublisher          *mevent.Publisher
	PlayerSummaryRepository lobby.PlayerSummaryReadRepository
	PlayerSummaryTranslator translate.LobbyPlayerSummaryTranslator
}

type ClientNotification interface {
	MarshallBinary() ([]byte, error)
}

func NewLobbyChannelEventNotifier(publisher *mevent.Publisher, summary lobby.PlayerSummaryReadRepository, translator translate.LobbyPlayerSummaryTranslator) *LobbyChannelEventNotifier {
	var subscriber = &LobbyChannelEventNotifier{EventPublisher: publisher, PlayerSummaryRepository: summary, PlayerSummaryTranslator: translator}
	publisher.Subscribe(subscriber, lobby.ParticipantCreatedEvent{}, lobby.ParticipantDeletedEvent{}, lobby.ParticipantReadyEvent{}, lobby.ParticipantUnreadyEvent{})
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
	}
}

func (s *LobbyChannelEventNotifier) HandleParticipantDeleted(event lobby.ParticipantDeletedEvent) error {

	var notification = &protomulti.ParticipantLeaveNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerId:   event.PlayerID().String(),
		PlayerSlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiNotificationType_PARTICIPANT_LEAVE, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantAdded(event lobby.ParticipantCreatedEvent) error {

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

	return s.publish(event.Context(), protomulti.MultiNotificationType_PARTICIPANT_JOIN, event.LobbyID(), notification)

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

	return s.publish(event.Context(), protomulti.MultiNotificationType_PARTICIPANT_DECK_CHANGE, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantReady(event lobby.ParticipantReadyEvent) error {

	var notification = &protomulti.ParticipantReadyNotification{
		LobbyId:    event.LobbyID().String(),
		PlayerSlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiNotificationType_PARTICIPANT_READY, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) HandleParticipantUnready(event lobby.ParticipantUnreadyEvent) error {

	var notification = &protomulti.ParticipantUnreadyNotification{
		LobbyId:   event.LobbyID().String(),
		PartySlot: int32(event.SlotIndex()),
	}

	return s.publish(event.Context(), protomulti.MultiNotificationType_PARTICIPANT_UNREADY, event.LobbyID(), notification)

}

func (s *LobbyChannelEventNotifier) publish(ctx context.Context, t protomulti.MultiNotificationType, id uuid.UUID, n ClientNotification) error {
	bytes, err := n.MarshallBinary()
	if err != nil {
		return err
	}
	var message = consumer.NewLobbyClientNotificationEvent(ctx, t, id, bytes)
	s.EventPublisher.Notify(message)
	return nil
}
