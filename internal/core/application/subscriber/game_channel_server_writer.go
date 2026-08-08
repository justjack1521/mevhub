package subscriber

import (
	"log/slog"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
	"time"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

type GameChannelServerWriter struct {
	Server                *server.GameServerHost
	InstanceRepository    port.GameInstanceRepository
	PartyRepository       port.GamePartyReadRepository
	ParticipantRepository port.GameParticipantReadRepository
	Logger                *slog.Logger
}

func NewGameChannelServerWriter(svr *server.GameServerHost, publisher *mevent.Publisher, instances port.GameInstanceRepository, party port.GamePartyReadRepository, participants port.GameParticipantReadRepository, logger *slog.Logger) *GameChannelServerWriter {
	var writer = &GameChannelServerWriter{Server: svr, InstanceRepository: instances, PartyRepository: party, ParticipantRepository: participants, Logger: logger}
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{}, game.PartyCreatedEvent{}, game.ParticipantCreatedEvent{}, session.InstanceDeletedEvent{}, game.PlayerDisconnectedEvent{}, game.PlayerReconnectedEvent{})
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
	case game.PlayerDisconnectedEvent:
		w.HandlePlayerDisconnected(actual)
	case game.PlayerReconnectedEvent:
		w.HandlePlayerReconnected(actual)
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
	if err := w.InstanceRepository.Delete(event.Context(), event.InstanceID()); err != nil {
		w.Logger.With("instance.id", event.InstanceID().String(), "error", err.Error()).Error("failed to delete game instance")
	}
}

func (w *GameChannelServerWriter) HandlePartyCreated(event game.PartyCreatedEvent) {
	party, err := w.PartyRepository.Query(event.Context(), event.GameID(), event.PartyIndex())
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

func (w *GameChannelServerWriter) HandlePlayerDisconnected(event game.PlayerDisconnectedEvent) {
	w.Server.ActionChannel <- &server.GameActionRequest{
		GameID: event.GameID(),
		Action: action.NewPlayerDisconnectAction(event.GameID(), event.PlayerID(), time.Now().UTC()),
	}
}

func (w *GameChannelServerWriter) HandlePlayerReconnected(event game.PlayerReconnectedEvent) {
	w.Server.ActionChannel <- &server.GameActionRequest{
		GameID: event.GameID(),
		Action: action.NewPlayerReconnectAction(event.PlayerID()),
	}
}
