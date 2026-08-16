package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
	"time"
)

// gameCatchUpTimeout bounds the wait for the game loop's reply. The request is
// answered from inside the loop, so a game that has already been unregistered
// (or a host under load) would otherwise leave a caller with no deadline of its
// own waiting forever.
const gameCatchUpTimeout = time.Second * 5

// GameCatchUpCommand asks the live game to resend the notifications a client
// missed. Unusually for a command it carries a result: the caller needs the
// range that was actually replayed, and the remaining turn time, which is the
// one piece of state the replayed events cannot express.
type GameCatchUpCommand struct {
	BasicCommand
	LastSequence uint64
	Result       action.CatchUpResult
}

func (c *GameCatchUpCommand) CommandName() string {
	return "game.catch_up"
}

func NewGameCatchUpCommand(last uint64) *GameCatchUpCommand {
	return &GameCatchUpCommand{LastSequence: last}
}

type GameCatchUpCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewGameCatchUpCommandHandler(sessions port.SessionInstanceReadRepository, svr *server.GameServerHost) *GameCatchUpCommandHandler {
	return &GameCatchUpCommandHandler{SessionRepository: sessions, GameServerHost: svr}
}

func (h *GameCatchUpCommandHandler) Handle(ctx Context, cmd *GameCatchUpCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	// The session is authoritative about which game the caller is in. A caller
	// not in one has nothing to catch up on; that is not an error.
	if uuid.Equal(current.GameID, uuid.Nil) {
		cmd.Result = emptyCatchUpResult(cmd.LastSequence)
		return nil
	}

	// Buffered by one, so the publisher's reply neither blocks nor is lost if
	// this handler has already given up on it.
	var result = make(chan action.CatchUpResult, 1)

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewCatchUpAction(current.GameID, current.PlayerID, cmd.LastSequence, result),
	}

	h.GameServerHost.ActionChannel <- request

	var timeout = time.NewTimer(gameCatchUpTimeout)
	defer timeout.Stop()

	select {
	case cmd.Result = <-result:
		return nil
	case <-timeout.C:
		cmd.Result = emptyCatchUpResult(cmd.LastSequence)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

// emptyCatchUpResult reports an inverted, and so deliberately empty, range: the
// client is told there is nothing to send rather than being handed a zeroed
// range it would read as the start of the game.
func emptyCatchUpResult(last uint64) action.CatchUpResult {
	return action.CatchUpResult{FromSequence: last + 1, ToSequence: last, TurnRemainingMs: -1}
}
