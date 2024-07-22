package subscriber

import (
	"github.com/justjack1521/mevium/pkg/mevent"
	"mevhub/internal/app/server"
	"mevhub/internal/domain/game"
)

type GameChannelServerWriter struct {
	Server                *server.GameServerHost
	EventPublisher        *mevent.Publisher
	InstanceRepository    game.InstanceReadRepository
	ParticipantRepository game.PlayerParticipantReadRepository
}

func NewGameChannelServerWriter(server *server.GameServerHost, publisher *mevent.Publisher, instances game.InstanceRepository, participants game.PlayerParticipantReadRepository) *GameChannelServerWriter {
	var writer = &GameChannelServerWriter{Server: server, EventPublisher: publisher, InstanceRepository: instances, ParticipantRepository: participants}
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{}, game.ParticipantCreatedEvent{})
	return writer
}

func (w *GameChannelServerWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		w.HandleInstanceCreated(actual)
	case game.InstanceDeletedEvent:
		w.HandleInstanceDelete(actual)
	case game.ParticipantCreatedEvent:
		w.HandleParticipantCreated(actual)
	}
}

func (w *GameChannelServerWriter) HandleInstanceCreated(event game.InstanceCreatedEvent) {
	instance, err := w.InstanceRepository.Get(event.Context(), event.InstanceID())
	if err != nil {
		return
	}
	w.Server.Register <- w.Server.NewLiveGameChannel(instance)
}

func (w *GameChannelServerWriter) HandleInstanceDelete(event game.InstanceDeletedEvent) {
	w.Server.Unregister <- event.InstanceID()
}

func (w *GameChannelServerWriter) HandleParticipantCreated(event game.ParticipantCreatedEvent) {

	participant, err := w.ParticipantRepository.Query(event.Context(), event.InstanceID(), event.PlayerSlot())
	if err != nil {
		return
	}

	w.Server.ActionChannel <- &server.GameActionRequest{
		InstanceID: event.InstanceID(),
		Action: &game.PlayerAddAction{
			UserID:    participant.UserID,
			PlayerID:  participant.PlayerID,
			PartySlot: participant.PlayerSlot,
		},
	}
}
