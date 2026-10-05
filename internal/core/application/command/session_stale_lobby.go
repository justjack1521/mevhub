package command

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/lobby"
	"mevhub/internal/core/domain/session"
	"mevhub/internal/core/port"
)

// releaseStaleLobby clears a session's lobby pointer when the lobby it names no longer
// exists. Lobby instances expire out of Redis without emitting anything, so a session can
// outlive its lobby and would otherwise refuse every lobby until the session itself expires.
// It returns true when the session is free to join a lobby.
func releaseStaleLobby(ctx context.Context, sessions port.SessionInstanceRepository, lobbies port.LobbyInstanceReadRepository, current *session.Instance) (bool, error) {

	if current.CanJoinLobby() {
		return true, nil
	}

	_, err := lobbies.QueryByID(ctx, current.LobbyID)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, lobby.ErrLobbyInstanceNotFound) {
		return false, err
	}

	current.LobbyID = uuid.Nil
	current.PartySlot = 0
	if err := sessions.Update(ctx, current); err != nil {
		return false, err
	}

	return true, nil

}
