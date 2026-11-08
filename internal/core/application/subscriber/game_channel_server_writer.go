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
	publisher.Subscribe(writer, game.InstanceCreatedEvent{}, game.InstanceDeletedEvent{}, session.InstanceDeletedEvent{}, game.PlayerDisconnectedEvent{}, game.PlayerReconnectedEvent{})
	return writer
}

func (w *GameChannelServerWriter) Notify(event mevent.Event) {
	switch actual := event.(type) {
	case game.InstanceCreatedEvent:
		w.HandleInstanceCreated(actual)
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

// HandleInstanceCreated builds the live game already holding its parties and
// players, then registers it. Built populated rather than filled in after
// registration: the parties and players used to be replayed as actions once the
// host had the game, a goroutine hop and several Redis reads behind, so a fast
// client's ready could reach the game before its own player did and be dropped.
// Now no action can reach the game before the players it is for.
//
// The read models are already written: GamePartyWriter and, through
// PartyCreated, GameParticipantWriter subscribe to InstanceCreated ahead of this
// writer and run to completion first, synchronously.
func (w *GameChannelServerWriter) HandleInstanceCreated(event game.InstanceCreatedEvent) {
	instance, err := w.InstanceRepository.Get(event.Context(), event.InstanceID())
	if err != nil {
		w.Logger.With("instance.id", event.InstanceID().String(), "error", err.Error()).Error("failed to get created game; it will not go live")
		return
	}
	var svr = w.Server.NewLiveGameChannel(instance)
	if err := w.seed(event, svr); err != nil {
		w.Logger.With("instance.id", event.InstanceID().String(), "error", err.Error()).Error("failed to query parties for created game; it will not go live")
		return
	}
	w.Server.Register <- svr
}

// seed adds the persisted parties and participants to the not-yet-started game.
// Only the party query failing is fatal - a game with no parties has nobody to
// play it. A party whose participants cannot be read, or a single add that
// fails, is logged and the rest still go in, as the old replay did.
func (w *GameChannelServerWriter) seed(event game.InstanceCreatedEvent, svr *server.GameServer) error {

	parties, err := w.PartyRepository.QueryAll(event.Context(), event.InstanceID())
	if err != nil {
		return err
	}

	for _, party := range parties {

		if err := svr.Seed(action.NewPartyAddAction(party.SysID, party.Index)); err != nil {
			w.Logger.With("instance.id", event.InstanceID().String(), "party.id", party.SysID.String(), "error", err.Error()).Error("failed to add party to created game")
			continue
		}

		participants, err := w.ParticipantRepository.QueryAll(event.Context(), party.SysID)
		if err != nil {
			w.Logger.With("instance.id", event.InstanceID().String(), "party.id", party.SysID.String(), "error", err.Error()).Error("failed to query participants for created game")
			continue
		}

		for _, participant := range participants {
			if err := svr.Seed(action.NewPlayerAddAction(participant.UserID, participant.PlayerID, party.SysID, participant.PlayerSlot)); err != nil {
				w.Logger.With("instance.id", event.InstanceID().String(), "party.id", party.SysID.String(), "player.id", participant.PlayerID.String(), "error", err.Error()).Error("failed to add player to created game")
			}
		}

	}

	return nil

}

func (w *GameChannelServerWriter) HandleInstanceDelete(event game.InstanceDeletedEvent) {
	w.Server.Unregister <- event.InstanceID()
	if err := w.InstanceRepository.Delete(event.Context(), event.InstanceID()); err != nil {
		w.Logger.With("instance.id", event.InstanceID().String(), "error", err.Error()).Error("failed to delete game instance")
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
		Action: action.NewPlayerReconnectAction(event.PlayerID(), time.Now().UTC()),
	}
}
