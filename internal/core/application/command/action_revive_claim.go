package command

import (
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/application/server"
	"mevhub/internal/core/domain/game/action"
	"mevhub/internal/core/port"
	"time"
)

// reviveClaimReplyTimeout bounds how long a claim waits for the game loop's
// answer. The loop drains actions immediately when live, so anything close to
// this deadline means the game is not hosted here, shut down mid-route, or
// badly stalled — all cases where the caller must be told "no" rather than
// left guessing whether their claim landed.
const reviveClaimReplyTimeout = time.Second * 3

// PlayerReviveClaimCommand asks for the decaying exclusive right to revive a
// dead player before the caller spends a revive item (BattleRevive on the game
// service). A nil target means claiming a self-revive. Unlike every other game
// command this one carries an answer back: Result is populated by the handler
// before Handle returns, and no reply in time is a denial — fail closed, so
// the client never consumes an item on an unconfirmed claim.
type PlayerReviveClaimCommand struct {
	BasicCommand
	TargetPlayerID uuid.UUID
	Result         action.PlayerReviveClaimResult
}

func (e *PlayerReviveClaimCommand) CommandName() string {
	return "action.revive.claim"
}

func NewPlayerReviveClaimCommand(target uuid.UUID) *PlayerReviveClaimCommand {
	return &PlayerReviveClaimCommand{TargetPlayerID: target}
}

type PlayerReviveClaimCommandHandler struct {
	SessionRepository port.SessionInstanceReadRepository
	GameServerHost    *server.GameServerHost
}

func NewPlayerReviveClaimCommandHandler(sessions port.SessionInstanceReadRepository, server *server.GameServerHost) *PlayerReviveClaimCommandHandler {
	return &PlayerReviveClaimCommandHandler{SessionRepository: sessions, GameServerHost: server}
}

func (h *PlayerReviveClaimCommandHandler) Handle(ctx Context, cmd *PlayerReviveClaimCommand) error {

	current, err := h.SessionRepository.QueryByID(ctx, ctx.UserID())
	if err != nil {
		return err
	}

	var target = cmd.TargetPlayerID
	if uuid.Equal(target, uuid.Nil) {
		target = current.PlayerID
	}

	// Buffered so the game loop's reply never blocks on us, even if we have
	// already timed out and stopped listening.
	var response = make(chan action.PlayerReviveClaimResult, 1)

	var request = &server.GameActionRequest{
		GameID:  current.GameID,
		PartyID: current.LobbyID,
		Action:  action.NewPlayerReviveClaimAction(current.GameID, target, current.PlayerID, time.Now().UTC(), response),
	}

	h.GameServerHost.ActionChannel <- request

	select {
	case result := <-response:
		cmd.Result = result
	case <-time.After(reviveClaimReplyTimeout):
		// An orphaned action (game not hosted, shut down mid-route) never
		// replies at all; this deadline is what turns that silence into a
		// denial instead of a hung RPC.
		cmd.Result = action.PlayerReviveClaimResult{Reason: action.ReviveClaimDenyReasonUnavailable}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil

}
