package lobby

import (
	"errors"
	"fmt"
)

var (
	// ErrLobbyInstanceNotFound is wrapped by repositories when the instance does not
	// exist. Lobby instances expire out of Redis without emitting anything, so callers
	// holding a lobby id from a session use this to tell "gone" from a real failure.
	ErrLobbyInstanceNotFound    = errors.New("lobby instance not found")
	ErrFailedQueryLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to query lobby instance: %w", err)
	}
)

var (
	ErrFailedCreateLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to create lobby instance: %w", err)
	}
	ErrFailedDeleteLobbyInstance = func(err error) error {
		return fmt.Errorf("failed to delete lobby instance: %w", err)
	}
)
