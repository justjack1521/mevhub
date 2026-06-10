package subscriber

import (
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type GameChannelServerWriter struct {
	Server                *server.GameServerHost
	EventPublisher        *mevent.Publisher
	InstanceRepository    port.GameInstanceReadRepository
	PartyRepository       port.GamePartyReadRepository
	ParticipantRepository port.GamePlayerReadRepository
}

func NewGameChannelServerWriter(server *server.GameServerHost, publisher *mevent.Publisher, instances port.GameInstanceRepository, party port.GamePartyReadRepository, participants port.GamePlayerReadRepository) *GameChannelServerWriter {
	var writer = &GameChannelServerWriter{Server: server, EventPublisher: publisher, InstanceRepository: instances, PartyRepository: party, ParticipantRepository: participants}
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{}, game.PartyCreatedEvent{}, game.ParticipantCreatedEvent{}, session.InstanceDeletedEvent{})
	return writer
}

func (w *GameChannelServerWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		w.HandleInstanceCreated(actual)
	case game.InstanceDeletedEvent:
		w.HandleInstanceDelete(actual)
	case game.PartyCreatedEvent:
		w.HandlePartyCreated(actual)
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
		Action:  action.NewPartyAddAction(party.SysID, party.Index),
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
		Action:  action.NewPlayerAddAction(participant.UserID, participant.PlayerID, event.PartyID(), participant.PlayerSlot),
	}
}

func (w *GameChannelServerWriter) HandleSessionDeleted(event session.InstanceDeletedEvent) {
	if uuid.Equal(event.GameID(), uuid.Nil) {
		return
	}
	w.Server.ActionChannel <- &server.GameActionRequest{
		GameID:  event.GameID(),
		PartyID: event.LobbyID(),
		Action:  action.NewPlayerRemoveAction(event.GameID(), event.LobbyID(), event.UserID(), event.PlayerID()),
	}
}
