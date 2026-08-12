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
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceRegisteredEvent{}, game.InstanceDeletedEvent{}, session.InstanceDeletedEvent{}, game.PlayerDisconnectedEvent{}, game.PlayerReconnectedEvent{})
	return writer
}

func (w *GameChannelServerWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		w.HandleInstanceCreated(actual)
	case game.InstanceRegisteredEvent:
		w.HandleInstanceRegistered(actual)
	case game.InstanceDeletedEvent:
		w.HandleInstanceDelete(actual)
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

// HandleInstanceRegistered replays the persisted parties and participants into
// the live game. Populating from the read models here rather than from the
// PartyCreated/ParticipantCreated events is what makes the ordering safe: those
// events fire while the server is still queued on Register, so their actions
// would arrive before the host has the game in its map and be dropped as
// orphaned. By registration time the read models are already written, because
// GamePartyWriter runs to completion on InstanceCreated before this writer
// queues the registration.
func (w *GameChannelServerWriter) HandleInstanceRegistered(event game.InstanceRegisteredEvent) {

	parties, err := w.PartyRepository.QueryAll(event.Context(), event.InstanceID())
	if err != nil {
		w.Logger.With("instance.id", event.InstanceID().String(), "error", err.Error()).Error("failed to query parties for registered game")
		return
	}

	for _, party := range parties {

		w.Server.ActionChannel <- &server.GameActionRequest{
			GameID:  event.InstanceID(),
			PartyID: party.SysID,
			Action:  action.NewPartyAddAction(party.SysID, party.Index),
		}

		participants, err := w.ParticipantRepository.QueryAll(event.Context(), party.SysID)
		if err != nil {
			w.Logger.With("instance.id", event.InstanceID().String(), "party.id", party.SysID.String(), "error", err.Error()).Error("failed to query participants for registered game")
			continue
		}

		for _, participant := range participants {
			w.Server.ActionChannel <- &server.GameActionRequest{
				GameID:  event.InstanceID(),
				PartyID: party.SysID,
				Action:  action.NewPlayerAddAction(participant.UserID, participant.PlayerID, party.SysID, participant.PlayerSlot),
			}
		}

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
