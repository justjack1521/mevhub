package subscriber

import (
	"context"
	"testing"

	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/domain/session"

	uuid "github.com/satori/go.uuid"
)

func TestGameChannelServerWriter_AbandonedGameRemovesPlayer(t *testing.T) {
	var host = &server.GameServerHost{ActionChannel: make(chan *server.GameActionRequest, 1)}
	var writer = &GameChannelServerWriter{Server: host}
	var user, player, game = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

	writer.Notify(session.NewGameAbandonedEvent(context.Background(), user, player, game))

	select {
	case request := <-host.ActionChannel:
		if !uuid.Equal(request.GameID, game) {
			t.Errorf("routed to game %s, want %s", request.GameID, game)
		}
		remove, ok := request.Action.(*action.PlayerRemoveAction)
		if !ok {
			t.Fatalf("action is %T, want *action.PlayerRemoveAction", request.Action)
		}
		if !uuid.Equal(remove.GameID, game) || !uuid.Equal(remove.UserID, user) || !uuid.Equal(remove.PlayerID, player) {
			t.Errorf("removes game %s user %s player %s, want %s %s %s", remove.GameID, remove.UserID, remove.PlayerID, game, user, player)
		}
	default:
		t.Fatal("no action sent to the game host")
	}
}
