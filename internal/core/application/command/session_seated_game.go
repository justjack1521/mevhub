package command

import (
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"

	"github.com/justjack1521/mevium/pkg/mevent"
	uuid "github.com/satori/go.uuid"
)

// leaveSeatedGame takes the player out of the live game their session still
// names before they create or join another lobby.
//
// A session's GameID is cleared whenever the player leaves a game's participant
// list (the game ending, a session end, a grace eviction), so one still set here
// is a game that holds a seat for someone who is no longer playing it: a join
// that failed on the client while the server went ahead. That game would keep
// waiting on them, and keep sending them its notifications - delivered per
// player, not per game - into whatever battle they play next. The client has no
// way back into it, so it is left rather than the new lobby refused.
//
// Published at once, not queued on the command: the old game is done with
// whether or not the new lobby goes on to succeed, and the session no longer
// records it once this returns.
func leaveSeatedGame(ctx Context, publisher *mevent.Publisher, sessions port.SessionInstanceWriteRepository, current *session.Instance) error {

	if uuid.Equal(current.GameID, uuid.Nil) {
		return nil
	}

	var abandoned = current.GameID
	current.GameID = uuid.Nil
	if err := sessions.Update(ctx, current); err != nil {
		current.GameID = abandoned
		return err
	}

	publisher.Notify(session.NewGameAbandonedEvent(ctx, current.UserID, current.PlayerID, abandoned))

	return nil

}
