package subscriber

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

type GameChannelServerWriter struct {
	Server                *server.GameServerHost
	EventPublisher        *mevent.Publisher
	InstanceRepository    port.GameInstanceReadRepository
	PartyRepository       port.GamePartyReadRepository
	ParticipantRepository port.GamePlayerReadRepository
}

func NewGameChannelServerWriter(server *server.GameServerHost, publisher *mevent.Publisher, instances port.GameInstanceRepository, participants port.GamePlayerReadRepository) *GameChannelServerWriter {
	var writer = &GameChannelServerWriter{Server: server, EventPublisher: publisher, InstanceRepository: instances, ParticipantRepository: participants}
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{}, game.ParticipantCreatedEvent{}, session.InstanceDeletedEvent{})
	return writer
}

func (w *GameChannelServerWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		w.HandleInstanceCreated(actual)
	case game.InstanceDeletedEvent:
		w.HandleInstanceDelete(actual)
	case game.PartyCreatedEvent:

	case game.ParticipantCreatedEvent:
		w.HandleParticipantCreated(actual)
	case session.InstanceDeletedEvent:
		w.HandleSessionDeleted(actual)
	}
}

func (w *GameChannelServerWriter) HandleInstanceCreated(event game.InstanceCreatedEvent) {
	instance, err := w.InstanceRepository.Get(event.Context(), event.InstanceID())
	if err != nil {
		return
	}
	w.Server.Register <- w.Server.NewLiveGameChannel(instance)
	w.EventPublisher.Notify(game.NewInstanceRegisteredEvent(event.Context(), event.InstanceID()))
}

func (w *GameChannelServerWriter) HandleInstanceDelete(event game.InstanceDeletedEvent) {
	w.Server.Unregister <- event.InstanceID()
}

func (w *GameChannelServerWriter) HandlePartyCreated(event game.PartyCreatedEvent) {
	party, err := w.PartyRepository.Query(event.Context(), event.PartyID(), event.PartyIndex())
	if err != nil {
		return
	}
	w.Server.ActionChannel <- &server.GameActionRequest{
		GameID:  event.GameID(),
		PartyID: event.PartyID(),
		Action: &game.PartyAddAction{
			PartyID:    party.SysID,
			PartyIndex: party.Index,
		},
	}
}

func (w *GameChannelServerWriter) HandleParticipantCreated(event game.ParticipantCreatedEvent) {

	participant, err := w.ParticipantRepository.Query(event.Context(), event.PartyID(), event.PlayerSlot())
	if err != nil {
		return
	}

	w.Server.ActionChannel <- &server.GameActionRequest{
		GameID:  event.GameID(),
		PartyID: event.PartyID(),
		Action: &game.PlayerAddAction{
			UserID:    participant.UserID,
			PlayerID:  participant.PlayerID,
			PartySlot: participant.PlayerSlot,
			PartyID:   event.PartyID(),
		},
	}
}

func (w *GameChannelServerWriter) HandleSessionDeleted(event session.InstanceDeletedEvent) {
	w.Server.ActionChannel <- &server.GameActionRequest{
		PartyID: event.LobbyID(),
		Action: &game.PlayerRemoveAction{
			InstanceID: event.LobbyID(),
			UserID:     event.UserID(),
			PlayerID:   event.PlayerID(),
		},
	}
}
