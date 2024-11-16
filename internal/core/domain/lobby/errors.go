package lobby

import (
	"fmt"
)

var (
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
